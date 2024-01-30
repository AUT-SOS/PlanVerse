import React, { useCallback, useState } from "react";
import { Text2 } from "../../../ui/Text";
import { State, Task } from "../../../utils/types";
import styles from "./Column.module.scss";
import { a, useTransition } from "@react-spring/web";
import { TaskFC } from "../Task/Task";
import { shallowEqual, useDispatch } from "react-redux";
import { ProjectActions } from "../../../redux/slices/project.slice";
import { useParams } from "react-router-dom";

type Props = {
  column: State;
};

function hexToRgb(hex: string) {
  var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result
    ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16),
      }
    : null;
}

export const Column: React.FC<Props> = (props) => {
  const color = hexToRgb(props.column.back_ground_color);

  const projId = useParams().id!;
  const newTask = {
    title: "New Task",
    id: projId,
    back_ground_color: "",
    description: "",
    performers: [],
  };

  const dispatch = useDispatch();

  /*const transition = useTransition(showNewTask, {
    from: {
      x: -100,
      opacity: "0",
    },
    enter: {
      x: 0,
      opacity: "1",
    },
  });*/

  const handleTaskAdd = useCallback(() => {
    dispatch(
      ProjectActions.createTask({
        ...newTask,
        state_id: props.column.state_id,
      })
    );
  }, [newTask]);

  return (
    <div
      className={styles.Column}
      style={{
        backgroundColor: `rgba(${color?.r},${color?.g},${color?.b}, 0.5) `,
      }}
    >
      <Text2 className={styles.ColumnTitle} text={props.column.title} />
      {props.column &&
        props.column.tasks &&
        props.column.tasks.map((item, index) => (
          <TaskFC task={item} key={index} state_id={props.column.state_id} />
        ))}
      <div className={styles.AddTask} onClick={handleTaskAdd}>
        + Add Task
      </div>
    </div>
  );
};

