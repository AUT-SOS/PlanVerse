import { useDispatch } from "react-redux";
import { SpinningLoading } from "../../ui/SpinningLoading";
import { State } from "../../utils/types";
import { Column } from "./Column/Column";
import styles from "./TaskBoard.module.scss";
import { ProjectActions } from "../../redux/slices/project.slice";

type Props = {
  states: State[];
  projectId: string;
};

export const TaskBoard: React.FC<Props> = (props) => {
  console.log(">>S", props.states);

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
      {props.states.map((item, index) => (
        <Column column={item} key={index} />
      ))}
      <div className={styles.AddState} onClick={addState}>
        + Add State
      </div>
    </div>
  ) : <SpinningLoading size={40}/>;
};
