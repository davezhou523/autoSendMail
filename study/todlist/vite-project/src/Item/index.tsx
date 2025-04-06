import {Component} from "react";
import './index.css'
export default class Item extends Component {
    render(){
        return (
            <>
                <li>
                    <label>
                        <input type="chekbox"/>
                        <span>xxx</span>
                    </label>
                    <button className="btn btn-danger" style={{display:'none'}}>删除</button>
                </li>
            </>
        )
    }
}