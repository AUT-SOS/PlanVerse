import React, { useState } from "react";
import { Task } from "../../../utils/types";
import styles from "./Task.module.scss";
import { Text1, Text2 } from "../../../ui/Text";
import { a } from "@react-spring/web";
import classNames from "classnames";
import { useDispatch } from "react-redux";
import { ProjectActions } from "../../../redux/slices/project.slice";
import { useParams } from "react-router-dom";

type Props = React.HTMLProps<HTMLDivElement> & {
  task: Task;
  state_id: string;
};

export const TaskFC: React.FC<Props> = (props) => {
  const [taskInfo, setTaskInfo] = useState<Task>(props.task);
  
  const projId = useParams().id!;

  const dispatch = useDispatch();
  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setTaskInfo((prev) => ({ ...prev, title: event.target.value }));
  };
  const handleBlur = () => {
    console.log(">>TT", {
      id: projId,  
      state_id: props.state_id,
      title: taskInfo.title,
      back_ground_color: taskInfo.back_ground_color,
      description: taskInfo.description,
      task_id: taskInfo.task_id,
    });
    dispatch(
      ProjectActions.editTask({
        id: projId,  
        state_id: props.state_id,
        title: taskInfo.title,
        back_ground_color: taskInfo.back_ground_color,
        description: taskInfo.description,
        task_id: taskInfo.task_id,
      })
    );
  };
  

  return (
    <div className={styles.Wrapper} draggable>
      <TextTitle
        onBlur={handleBlur}
        value={taskInfo.title}
        onChange={handleEditName}
      />
    </div>
  );
};

type TitleProps = React.HTMLProps<HTMLInputElement>;

const TextTitle: React.FC<TitleProps> = (props) => {
  return (
    <input
      type="text"
      className={classNames(styles.TaskTitle, props.className)}
      {...props}
      value={props.value}
    />
  );
};
