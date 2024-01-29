import { useNavigate, useParams } from "react-router-dom";
import styles from "./Board.module.scss";
import { project1 } from "../../utils/testCase";
import { Title } from "../../ui/Title";
import React, {
  PropsWithChildren,
  useCallback,
  useEffect,
  useState,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { SpinningLoading } from "../../ui/SpinningLoading";
import { ProjectActions } from "../../redux/slices/project.slice";
import { Text0, Text1 } from "../../ui/Text";
import { Home } from "../../ui/Icons/Home";
import { Members } from "../../ui/Icons/Members";
import { Settings } from "../../ui/Icons/Settings";
import { a, useTransition } from "@react-spring/web";
import classNames from "classnames";
import { BackIcon } from "../../ui/Icons/Back";
import { Member } from "../../utils/types";
import { MemberItem } from "../../ui/MemberItem";

type Props = {};

enum SliderTypes {
  Settings,
  Members,
}

export const Board: React.FC<Props> = (props) => {
  const projId = useParams().id;
  const navigate = useNavigate();
  if (!projId) {
    navigate("/home");
    return;
  }

  const transition = useTransition(projId, {
    from: {
      y: 100,
      opacity: "0",
    },
    enter: {
      y: 0,
      opacity: "1",
    },
    leave: {
      y: -200,
      opacity: "0",
    },
  });

  const dispatch = useDispatch();
  const { project, members } = useSelector((state: RootState) => ({
    project: state.project.fullProject,
    members: state.project.members,
  }));

  useEffect(() => {
    dispatch(ProjectActions.getFullProject(projId));
  }, []);
  const [visible, setVisible] = useState(false);
  const [sliderContent, setSliderContent] = useState(SliderTypes.Members);

  const handleSettings = useCallback(() => {
    setSliderContent(SliderTypes.Settings);
    setVisible(true);
  }, []);

  const handleMembers = useCallback(() => {
    setSliderContent(SliderTypes.Members);
    setVisible(true);
  }, []);

  const sliderTitle =
    sliderContent === SliderTypes.Members ? "Members" : "Project Setting";

  return project && members ? (
    <div className={styles.BoardWrapper}>
      <MoreInfoSlider
        title={sliderTitle}
        visible={visible}
        setVisible={setVisible}
      >
        {sliderContent === SliderTypes.Members ? (
          <ProjectMembers
            members={members}
            ownerId={project.owner_id}
            projectId={project.project_id}
          />
        ) : (
          <div></div>
        )}
      </MoreInfoSlider>
      <img className={styles.ProjImg} src={project.picture} alt="" />
      {transition((style, state) => (
        <a.div style={style} className={styles.NavBarWrapper}>
          <Text0 text={project.title} className={styles.ProjectName} />
          <div className={styles.Navbar}>
            <Members
              onClick={handleMembers}
              size={30}
              color={"white"}
              className={styles.Icon}
            />
            <Settings
              onClick={handleSettings}
              size={30}
              color={"white"}
              className={styles.Icon}
            />
            <Home
              size={30}
              color={"white"}
              className={styles.Icon}
              onClick={() => navigate("/home")}
            />
          </div>
        </a.div>
      ))}

      <div className={styles.BoardContent}></div>
    </div>
  ) : (
    <div
      style={{
        width: "100vw",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <SpinningLoading size={50} />
    </div>
  );
};

type SliderProps = PropsWithChildren & {
  visible?: boolean;
  setVisible: (val: boolean) => void;
  title: string;
};

const MoreInfoSlider: React.FC<SliderProps> = (props) => {
  return (
    <div
      className={classNames(styles.InfoSlider, {
        [styles.visible]: props.visible,
      })}
    >
      <div className={styles.SliderHeader}>
        <Text1 text={props.title} />
        <BackIcon
          size={30}
          color={"var(--color-neutrals-n-600)"}
          className={styles.BackIcon}
          onClick={() => props.setVisible(false)}
        />
      </div>
      {props.children}
    </div>
  );
};

const ProjectMembers: React.FC<{
  members: Member[];
  ownerId: string;
  projectId: string;
}> = (props) => {
  const myId = useSelector((state: RootState) => state.auth.myId);
  const amIAdmin = props.members.find((item) => item.id === myId)?.is_admin;

  return (
    <div className={styles.MemberList}>
      {props.members.map((item, index) => {
        return (
          <MemberItem
            amIAdmin={amIAdmin}
            ownerId={props.ownerId}
            amIOwner={props.ownerId === myId}
            member={item}
            key={index}
            projectId={props.projectId}
          />
        );
      })}
    </div>
  );
};
