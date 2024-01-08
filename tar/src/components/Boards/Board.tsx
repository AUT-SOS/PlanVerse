import { useParams } from "react-router-dom";
import styles from "./Board.module.scss";
import { project1 } from "../../utils/testCase";
import { Title } from "../../ui/Title";

type Props = {};

export const Board: React.FC<Props> = (props) => {
  const groupId = useParams();
  const project = project1;
  console.log(">>", groupId);
  return (
    <div
      className={styles.BoardWrapper}
      style={{ background: `url(${project.background})` }}
    >
      <div className={styles.BackgroundFilter}>
            <Title text={project.name} className={styles.ProjectName}/>
            
      </div>
    </div>
  );
};
