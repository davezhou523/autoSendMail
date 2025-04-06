import {Component} from "react";
import './index.css'
import {nanoid} from "nanoid";
export default class Header extends Component {
    handleKeyDown=(event)=>{
        console.log(event.target.value);
        if (event.keyCode !=13){
            return;
        }
        if (event.target.value.trim()===""){
            alert("not empty")
            return;
        }
        const todoObj={id:nanoid(),name:event.target.value,done:false}
        // @ts-ignore
        //子传父参数，父传参数为函数
        this.props.addTodo(todoObj)
    }
    render(){
        return (
            <>
                <div className="todo-header">
                    <input onKeyDown={this.handleKeyDown} type="text" placeholder="请输入任务名称，按回车确认"/>
                </div>
            </>
        )
    }
}