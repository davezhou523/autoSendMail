import { useState } from 'react'
import './App.css'
import Header from "./Header/index.jsx";
import List from "./List/index.jsx";
import Footer from "./Footer/index.jsx";

function App() {
    const [todos, setTodos] = useState([
        {id:1,name:"吃饭",done:true},
        {id:2,name:"睡觉",done:true},
        {id:3,name:"打代码",done:false},
        {id:4,name:"逛街",done:true},
    ]);

    const addTodo = (todoObj) => {
        setTodos(prev => [todoObj, ...prev]);
    }

    return (
        <div className="todo-container">
            <div className="todo-wrap">
                <Header addTodo={addTodo} />
                <List todos={todos} />
                <Footer />
            </div>
        </div>
    )
}

export default App