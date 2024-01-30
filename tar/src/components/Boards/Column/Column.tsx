import React, { useCallback, useState } from "react";
import { Text2 } from "../../../ui/Text";
import { State, Task } from "../../../utils/types";
import styles from "./Column.module.scss";
import { a, useTransition } from "@react-spring/web";
import { TaskFC, TextTitle } from "../Task/Task";
import { shallowEqual, useDispatch } from "react-redux";
import { ProjectActions } from "../../../redux/slices/project.slice";
import { useParams } from "react-router-dom";
import { SketchPicker } from "react-color";

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
  const [columnInfo, setColumnInfo] = useState<State>(props.column);
  const color = hexToRgb(columnInfo.back_ground_color);

  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setColumnInfo((prev) => ({ ...prev, title: event.target.value }));
  };

  const [display, setDisplay] = useState(false);
  const [sketchDisplay, setSketchDisplay] = useState(false);

  const handleBlur = () => {
    dispatch(
      ProjectActions.editState({
        project_id: projId,
        state_id: props.column.state_id,
        title: columnInfo.title,
        back_ground_color: columnInfo.back_ground_color,
        admin_access: columnInfo.admin_access,
        tasks: columnInfo.tasks,
      })
    );
  };

  const projId = useParams().id!;
  const newTask = {
    title: "New Task",
    project_id: projId,
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
      {sketchDisplay && (
        <SketchPicker
          color={columnInfo.back_ground_color}
          onChange={(e) =>
            setColumnInfo((prev) => ({ ...prev, back_ground_color: e.hex }))
          }
        />
      )}
      <div
        onMouseOver={() => setDisplay(true)}
        onMouseLeave={() => setDisplay(false)}
      >
        <TextTitle
          className={styles.ColumnTitle}
          value={columnInfo.title}
          onChange={handleEditName}
          onBlur={handleBlur}
        />
        {display && <div>
          
          </div>}
      </div>
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
