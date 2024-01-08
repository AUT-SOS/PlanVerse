import React from "react";
import { Task } from "../../../utils/types";
import "./Task.scss";

type Props =  React.HTMLProps<HTMLDivElement> & {
      task: Task;
}

const Task : React.FC<Props> = (props) => {
      return <div className={"TaskWrapper"}>

      </div>
}