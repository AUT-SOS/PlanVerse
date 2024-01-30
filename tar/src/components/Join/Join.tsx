import React, { useCallback, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ReqButton } from "../../ui/ReqButton";
import classNames from "classnames";
import { Members } from "../../ui/Icons/Members";
import "./Join.module.scss";
import {
  AuthState,
  Project,
  RequestState,
  RequestTypes,
} from "../../utils/types";
import { project1 } from "../../utils/testCase";
import { HollowButton } from "../../ui/HollowButton";
import { Title } from "../../ui/Title";
import { Background } from "../../ui/BackGround";
import styles from "./Join.module.scss";
import strings from "../../utils/text";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { ProjectActions } from "../../redux/slices/project.slice";
import { SpinningLoading } from "../../ui/SpinningLoading";
import { ReqActions } from "../../redux/slices/req.slice";
import { navigate } from "../../utils/configs";
import { useRequestStates } from "../../utils/hooks";
import { a, useTransition } from "@react-spring/web";

export const Join: React.FC = (props) => {
  const dispatch = useDispatch();
  const {isPending} = useRequestStates(RequestTypes.JoinProject)

  useEffect(() => {
    dispatch(ProjectActions.showProject(location.toString()));
  }, []);
  const { project } = useSelector((state: RootState) => ({
    project: state.project.joinProject,
  }));

  const transition = useTransition(project, {
    from: {
      y: 200,
      opacity: "0",
    },
    enter: {
      y: 0,
      opacity: "1",
    },
  });

  const handleJoin = useCallback(() => {
    if (!project) return;
    dispatch(
      ReqActions.setState({
        requestState: RequestState.Pending,
        reqType: RequestTypes.JoinProject,
      })
    );
    dispatch(ProjectActions.joinProject(project?.project_id));
  }, [project]);

  return project ? (
    <>
      <Title text={strings.palverse} href="/" />
      <Background className={styles.JoinWrapper}>
        {transition((style) => {
        return <a.div style={style} className={classNames(styles.JoinCard)}>
          <div className={styles.contentWrapper}>
            <img className={styles.GroupIMG} src={project.picture} alt="" />
            <div className={styles.GroupName}>{project.title}</div>
            <div className={styles.GroupMemberWrapper}>
              <div className={styles.GroupMembers}>
                {project.members.map((item, index) => {
                  return (
                    <img
                      key={index}
                      title={item.Username}
                      src={item.ProfilePic}
                      className={styles.MemberPrev}
                    />
                  );
                })}
              </div>
              <Members size={22} color={"var(--color-neutrals-n-500)"} />
              <p className={styles.MembersCount}>{project.members_number}</p>
            </div>
            <ReqButton
              onClick={handleJoin}
              isPending={isPending}
              text={strings.join.joinProject}
              style={{ width: "100%" }}
            />
          </div>
        </a.div>
        })}
      </Background>
    </>
  ) : (
    <SpinningLoading size={60} />
  );
};
