import { useDispatch } from "react-redux";
import { SpinningLoading } from "../../ui/SpinningLoading";
import { State, Task } from "../../utils/types";
import { Column } from "./Column/Column";
import styles from "./TaskBoard.module.scss";
import { ProjectActions } from "../../redux/slices/project.slice";

type Props = {
  states: State[];
  projectId: string;
  openTask: (task: Task, state_id: string) => void
};

export const TaskBoard: React.FC<Props> = (props) => {

  const dispatch = useDispatch();

  const addState = () => {
    dispatch(ProjectActions.createState({
      project_id: props.projectId,
      title: "New State",
      back_ground_color: "#444444",
      admin_access: false,
    }))
  }
  

  return props.states ? (
    <div className={styles.TaskBoard}>
      {props.states.map((item) => (
        <Column openTask={props.openTask} column={item} key={item.state_id} />
      ))}
      <div className={styles.AddState} onClick={addState}>
        + Add State
      </div>
    </div>
  ) : <SpinningLoading size={40}/>;
};
