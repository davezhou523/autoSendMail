import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Header from "./Header";
import List from "./List";
import Footer from "./Footer";

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
    <div className="todo-container">
    <div className="">
        <Header></Header>
        <List></List>
        <Footer></Footer>
    </div>
    </div>
    </>
  )
}

export default App
