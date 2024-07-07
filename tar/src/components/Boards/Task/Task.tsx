import React, { useEffect, useState } from "react";
import { Task } from "../../../utils/types";
import styles from "./Task.module.scss";
import { Text1, Text2 } from "../../../ui/Text";
import { a } from "@react-spring/web";
import classNames from "classnames";
import { shallowEqual, useDispatch, useSelector } from "react-redux";
import { ProjectActions } from "../../../redux/slices/project.slice";
import { useParams } from "react-router-dom";
import { RootState } from "../../../redux/store";

type Props = React.HTMLProps<HTMLDivElement> & {
  task: Task;
  state_id: string;
  openTask: (task: Task, state_id: string) => void;
};

export const TaskFC: React.FC<Props> = React.memo((props) => {
  const [taskInfo, setTaskInfo] = useState<Task>(props.task);
  const members = useSelector((state: RootState) => state.project.members);
  useEffect(() => {
    setTaskInfo(props.task);
  }, [props.task]);

  const projId = useParams().id!;

  const dispatch = useDispatch();
  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    event.preventDefault();
    setTaskInfo((prev) => ({ ...prev, title: event.target.value }));
    dispatch(
      ProjectActions.editTask({
        project_id: projId,
        state_id: props.state_id,
        title: event.target.value,
        back_ground_color: taskInfo.back_ground_color,
        description: taskInfo.description,
        task_id: taskInfo.task_id,
        index: taskInfo.index,
      })
    );
  };

  const handleClick = () => {
    props.openTask(props.task, props.state_id);
  };

  return (
    <div
      onClick={handleClick}
      className={styles.Wrapper}
      draggable
      onDragStart={(e) =>
        e.dataTransfer.setData("Task", JSON.stringify(props.task))
      }
    >
      <TextTitle
        onFocus={(e) => e.stopPropagation()}
        onClick={(e) => e.stopPropagation()}
        value={taskInfo.title}
        onChange={handleEditName}
      />
      {members && props.task.performers && (
        <div className={styles.AssignedList}>
          {props.task.performers.map((item) => {
            const target = members.find((a) => Number(a.id) == item);
            return (
              <img
                title={target?.username}
                className={styles.AssignedPic}
                src={
                  target?.profile_pic && target?.profile_pic.length > 0
                    ? target.profile_pic
                    : "/public/DefaultPFP.jpg"
                }
                alt=""
              />
            );
          })}
        </div>
      )}
    </div>
  );
}, shallowEqual);

type TitleProps = React.HTMLProps<HTMLInputElement>;

export const TextTitle: React.FC<TitleProps> = (props) => {
  return (
    <input
      type="text"
      {...props}
      className={classNames(styles.TaskTitle, props.className)}
      value={props.value}
    />
  );
};
