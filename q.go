package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Answer struct {
	QuestionID int
	AnswerID   int
}

type Cell struct {
	ID              int
	Mapping         *QuestionMapping
	Answers         []Answer
	EligibleAnswers []Answer
}

type QuestionMapping struct {
	ID                 int
	Country            string
	Cells              []*Cell
	AnswerIDToCells    map[int][]*Cell
	AnswerInMultiCells map[int]bool
}

type DataStore struct {
	CountryMappings    map[string][]*QuestionMapping
	CountryAnswerIndex map[string]map[int][]*QuestionMapping
}

var currentData atomic.Value

func main() {
	dataStore, err := loadData()
	if err != nil {
		log.Fatal("Failed to load data:", err)
	}
	currentData.Store(dataStore)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			newData, err := loadData()
			if err != nil {
				log.Println("Failed to refresh data:", err)
				continue
			}
			currentData.Store(newData)
			log.Println("Data refreshed successfully")
		}
	}()

	http.HandleFunc("/process", processHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadData() (*DataStore, error) {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	mappings, err := loadQuestionMappings(db)
	if err != nil {
		return nil, err
	}

	cells, err := loadCells(db, mappings)
	if err != nil {
		return nil, err
	}

	err = loadAnswers(db, cells)
	if err != nil {
		return nil, err
	}

	dataStore := &DataStore{
		CountryMappings:    make(map[string][]*QuestionMapping),
		CountryAnswerIndex: make(map[string]map[int][]*QuestionMapping),
	}

	for _, mapping := range mappings {
		dataStore.CountryMappings[mapping.Country] = append(dataStore.CountryMappings[mapping.Country], mapping)
		mapping.AnswerIDToCells = make(map[int][]*Cell)
		for _, cell := range mapping.Cells {
			for _, answer := range cell.Answers {
				mapping.AnswerIDToCells[answer.AnswerID] = append(mapping.AnswerIDToCells[answer.AnswerID], cell)
			}
		}
		mapping.AnswerInMultiCells = make(map[int]bool)
		for answerID, cells := range mapping.AnswerIDToCells {
			mapping.AnswerInMultiCells[answerID] = len(cells) > 1
		}
	}

	for _, cell := range cells {
		questionCounts := make(map[int]int)
		for _, answer := range cell.Answers {
			questionCounts[answer.QuestionID]++
		}
		var eligible []Answer
		for _, answer := range cell.Answers {
			if questionCounts[answer.QuestionID] == 1 {
				eligible = append(eligible, answer)
			}
		}
		cell.EligibleAnswers = eligible
	}

	for country, mappings := range dataStore.CountryMappings {
		dataStore.CountryAnswerIndex[country] = make(map[int][]*QuestionMapping)
		for _, mapping := range mappings {
			for answerID := range mapping.AnswerIDToCells {
				if !mapping.AnswerInMultiCells[answerID] {
					dataStore.CountryAnswerIndex[country][answerID] = append(dataStore.CountryAnswerIndex[country][answerID], mapping)
				}
			}
		}
	}

	return dataStore, nil
}

func loadQuestionMappings(db *sql.DB) ([]*QuestionMapping, error) {
	rows, err := db.Query("SELECT id, country FROM question_mappings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []*QuestionMapping
	for rows.Next() {
		var m QuestionMapping
		if err := rows.Scan(&m.ID, &m.Country); err != nil {
			return nil, err
		}
		mappings = append(mappings, &m)
	}
	return mappings, nil
}

func loadCells(db *sql.DB, mappings []*QuestionMapping) ([]*Cell, error) {
	rows, err := db.Query("SELECT id, question_mapping_id FROM question_mapping_cells")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mappingMap := make(map[int]*QuestionMapping)
	for _, m := range mappings {
		mappingMap[m.ID] = m
	}

	var cells []*Cell
	for rows.Next() {
		var cID, mID int
		if err := rows.Scan(&cID, &mID); err != nil {
			return nil, err
		}
		m := mappingMap[mID]
		cell := &Cell{ID: cID, Mapping: m}
		cells = append(cells, cell)
		m.Cells = append(m.Cells, cell)
	}
	return cells, nil
}

func loadAnswers(db *sql.DB, cells []*Cell) error {
	rows, err := db.Query("SELECT question_mapping_cell_id, question_id, answer_id FROM question_mapping_cell_answers")
	if err != nil {
		return err
	}
	defer rows.Close()

	cellMap := make(map[int]*Cell)
	for _, c := range cells {
		cellMap[c.ID] = c
	}

	for rows.Next() {
		var cellID, qID, aID int
		if err := rows.Scan(&cellID, &qID, &aID); err != nil {
			return err
		}
		cell := cellMap[cellID]
		if cell != nil {
			cell.Answers = append(cell.Answers, Answer{QuestionID: qID, AnswerID: aID})
		}
	}
	return nil
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Country string `json:"country"`
		Data    []struct {
			QuestionID int `json:"question_id"`
			AnswerID   int `json:"answer_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	dataStore := currentData.Load().(*DataStore)
	inputSet := make(map[Answer]bool)
	for _, item := range input.Data {
		inputSet[Answer{item.QuestionID, item.AnswerID}] = true
	}

	generatedSet := make(map[Answer]bool)
	countryIndex := dataStore.CountryAnswerIndex[input.Country]
	if countryIndex != nil {
		for _, item := range input.Data {
			answerID := item.AnswerID
			mappings := countryIndex[answerID]
			for _, mapping := range mappings {
				cells := mapping.AnswerIDToCells[answerID]
				if len(cells) == 0 {
					continue
				}
				cell := cells[0]
				for _, eligible := range cell.EligibleAnswers {
					if !inputSet[eligible] {
						generatedSet[eligible] = true
					}
				}
			}
		}
	}

	outputData := make([]Answer, 0, len(inputSet)+len(generatedSet))
	for a := range inputSet {
		outputData = append(outputData, a)
	}
	for a := range generatedSet {
		outputData = append(outputData, a)
	}

	sort.Slice(outputData, func(i, j int) bool {
		if outputData[i].QuestionID == outputData[j].QuestionID {
			return outputData[i].AnswerID < outputData[j].AnswerID
		}
		return outputData[i].QuestionID < outputData[j].QuestionID
	})

	response := struct {
		Data []struct {
			QuestionID int `json:"question_id"`
			AnswerID   int `json:"answer_id"`
		} `json:"data"`
	}{
		Data: make([]struct {
			QuestionID int `json:"question_id"`
			AnswerID   int `json:"answer_id"`
		}, len(outputData)),
	}

	for i, a := range outputData {
		response.Data[i].QuestionID = a.QuestionID
		response.Data[i].AnswerID = a.AnswerID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
