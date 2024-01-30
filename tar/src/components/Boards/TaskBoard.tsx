import { SpinningLoading } from "../../ui/SpinningLoading";
import { State } from "../../utils/types";
import { Column } from "./Column/Column";
import styles from "./TaskBoard.module.scss";

type Props = {
  states: State[];
  projectId: string;
};

export const TaskBoard: React.FC<Props> = (props) => {
  console.log(">>S", props.states);
  

  return props.states ? (
    <div className={styles.TaskBoard}>
      {props.states.map((item, index) => (
        <Column column={item} key={index} />
      ))}
    </div>
  ) : <SpinningLoading size={40}/>;
};
