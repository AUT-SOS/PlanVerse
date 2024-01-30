import React, { useCallback, useState } from "react";
import { State, Task } from "../../../utils/types";
import styles from "./Column.module.scss";
import { TaskFC, TextTitle } from "../Task/Task";
import { useDispatch, useSelector } from "react-redux";
import { ProjectActions } from "../../../redux/slices/project.slice";
import { useParams } from "react-router-dom";
import { SketchPicker } from "react-color";
import { ColorPicker } from "../../../ui/Icons/ColorPicker";
import { Delete } from "../../../ui/Icons/Delete";
import { a, useTransition } from "@react-spring/web";
import { RootState } from "../../../redux/store";

type Props = {
  column: State;
  openTask: (task: Task, state_id: string) => void

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
  const {members, myId} = useSelector((state: RootState) => ({members: state.project.members, myId: state.auth.myId}));
  const amIAdmin = members?.find((item) => item.id == myId)?.is_admin ?? false;


  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (!amIAdmin) return 
    setColumnInfo((prev) => ({ ...prev, title: event.target.value }));
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

  const [display, setDisplay] = useState(false);
  const [sketchDisplay, setSketchDisplay] = useState(false);

  const projId = useParams().id!;
  const newTask = {
    title: "New Task",
    project_id: projId,
    back_ground_color: "",
    description: "",
    performers: [],
    index: (props.column.tasks ?? []).length,
  };

  const dispatch = useDispatch();

  const transition = useTransition(display && amIAdmin, {
    from: {
      y: -50,
      opacity: "0",
    },
    enter: {
      y: 0,
      opacity: "1",
    },
    leave:{
      opacity: "0",
    }
    
  });

  const handleTaskAdd = useCallback(() => {
    dispatch(
      ProjectActions.createTask({
        ...newTask,
        state_id: props.column.state_id,
      })
    );
  }, [newTask]);

  const handleDeleteState = useCallback(() => {
    dispatch(
      ProjectActions.deleteState({
        projId,
        stateId: props.column.state_id,
      })
    );
  }, [projId, props.column]);

  const sortedTasks = props.column.tasks ? [...props.column.tasks] : [];
  sortedTasks.sort((a, b) => a.index - b.index);

  const handleDrop = (e: React.DragEvent) => {
    const task = JSON.parse(e.dataTransfer?.getData("Task") as string) as Task;
    
    dispatch(ProjectActions.changeState({
      project_id: projId,
      state_id: props.column.state_id,
      task_id: task.task_id
    }))
  }

  return (
    <div
    onDrop={handleDrop}
    onDragOver={(e) => e.preventDefault()}
      className={styles.Column}
      style={{
        backgroundColor: `rgba(${color?.r},${color?.g},${color?.b}, 0.5) `,
      }}
    >
      <div
        onMouseOver={() => setDisplay(true)}
        onMouseLeave={() => setDisplay(false)}
        className={styles.ColumnInf}
      >
        <TextTitle
          className={styles.ColumnTitle}
          value={columnInfo.title}
          onChange={handleEditName}
        />
        {transition((style, state) => state && (
          <a.div style={style} className={styles.EditWrapper}>
            <ColorPicker
              style={{ cursor: "pointer" }}
              color="var(--color-neutrals-n-00)"
              size={23}
              onClick={() => setSketchDisplay((prev) => !prev)}
            />
            <Delete
              style={{ cursor: "pointer" }}
              color="var(--color-neutrals-n-00)"
              size={23}
              onClick={handleDeleteState}
            />
          </a.div>
        ))}
        {sketchDisplay && (
          <SketchPicker
            color={columnInfo.back_ground_color}
            onChange={(e) => {
              setColumnInfo((prev) => ({ ...prev, back_ground_color: e.hex }));
              dispatch(ProjectActions.editState({project_id: projId, ...columnInfo}))
            }}
            className={styles.ColorPicker}
          />
        )}
      </div>
      {sortedTasks.map((item) => (
        <TaskFC openTask={props.openTask} task={item} key={item.task_id} state_id={props.column.state_id} />
      ))}
      <div className={styles.AddTask} onClick={handleTaskAdd}>
        + Add Task
      </div>
    </div>
  );
};
