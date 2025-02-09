package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Define Data Structures
type InputData struct {
	Country string `json:"country"`
	Data    []QA   `json:"data"`
}

type QA struct {
	QuestionID int `json:"question_id"`
	AnswerID   int `json:"answer_id"`
}

type QuestionMappingCell struct {
	CellID  int
	Answers map[int]int
}

type Service struct {
	db           *sql.DB
	dataMutex    sync.RWMutex
	questionData map[string]map[int][]QuestionMappingCell
}

func NewService(db *sql.DB) *Service {
	s := &Service{
		db:           db,
		questionData: make(map[string]map[int][]QuestionMappingCell),
	}
	go s.refreshData()
	return s
}

func (s *Service) refreshData() {
	for {
		s.loadData()
		time.Sleep(5 * time.Minute)
	}
}

// Load Data from MySQL into Memory
func (s *Service) loadData() {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return
	}
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()
	defer tx.Commit()

	rows, err := tx.Query("SELECT qmc.id, qm.country, qmc.id, qmc.id FROM question_mappings qm JOIN question_mapping_cells qmc ON qm.id = qmc.question_mapping_id")
	if err != nil {
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cellID, mappingID int
		var country string
		if err := rows.Scan(&mappingID, &country, &cellID); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		if _, exists := s.questionData[country]; !exists {
			s.questionData[country] = make(map[int][]QuestionMappingCell)
		}

		cell := QuestionMappingCell{CellID: cellID, Answers: make(map[int]int)}
		s.questionData[country][mappingID] = append(s.questionData[country][mappingID], cell)
	}

	rows, err = tx.Query("SELECT question_mapping_cell_id, question_id, answer_id FROM question_mapping_cell_answers")
	if err != nil {
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cellID, questionID, answerID int
		if err := rows.Scan(&cellID, &questionID, &answerID); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		for country, mappings := range s.questionData {
			for mappingID, cells := range mappings {
				for i, cell := range cells {
					if cell.CellID == cellID {
						s.questionData[country][mappingID][i].Answers[questionID] = answerID
					}
				}
			}
		}
	}
}

// Process API Requests
func (s *Service) processRequest(w http.ResponseWriter, r *http.Request) {
	var input InputData
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()

	responseData := make(map[int]int)
	for _, qa := range input.Data {
		responseData[qa.QuestionID] = qa.AnswerID
		for _, mapping := range s.questionData[input.Country] {
			for _, cell := range mapping {
				if _, exists := cell.Answers[qa.QuestionID]; exists {
					for qid, aid := range cell.Answers {
						if _, alreadyExists := responseData[qid]; !alreadyExists {
							responseData[qid] = aid
						}
					}
				}
			}
		}
	}

	output := struct {
		Data []QA `json:"data"`
	}{}
	for qid, aid := range responseData {
		output.Data = append(output.Data, QA{QuestionID: qid, AnswerID: aid})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func main() {
	db, err := sql.Open("mysql", "root:123abc@tcp(localhost:3306)/question")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	service := NewService(db)
	http.HandleFunc("/process", service.processRequest)
	log.Println("Service started on :8080")
	http.ListenAndServe(":8080", nil)
}
