import React, { useCallback, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { SpinningLoading } from "../../ui/SpinningLoading";
import {
  AuthState,
  Breakpoints,
  Project,
  RequestState,
  RequestTypes,
  SmallProject,
  User,
} from "../../utils/types";
import { useNavigate } from "react-router-dom";
import { Background } from "../../ui/BackGround";
import styles from "./Home.module.scss";
import { a, useTransition } from "@react-spring/web";
import { Text1, Text2, Text3 } from "../../ui/Text";
import { HollowButton, ReqButton1 } from "../../ui/ReqButton";
import classNames from "classnames";
import { Plus } from "../../ui/Icons/Plus";
import { useBreakPoints, useRequestStates } from "../../utils/hooks";
import { sProj1 } from "../../utils/testCase";
import { Members } from "../../ui/Icons/Members";
import { Logout } from "../../ui/Icons/Logout";
import { TextAreaInputBar, UsernameInputBar } from "../../ui/InputBar";
import { ReqActions } from "../../redux/slices/req.slice";
import { ProjectActions } from "../../redux/slices/project.slice";

enum HomeTypes {
  Overview,
  CreateProject,
  EditProfile,
}

export const Home: React.FC = (props) => {
  const { userInf, isAuth } = useSelector((state: RootState) => ({
    userInf: state.users.me,
    isAuth: state.auth.authState === AuthState.Authenticated,
  }));
  const [currState, setCurrState] = useState<HomeTypes>(HomeTypes.Overview);
  const navigate = useNavigate();
  if (!isAuth) {
    navigate("/login");
  }
  const transition = useTransition(currState, {
    from: {
      y: 200,
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

  const clearCookie = () => {
    document.cookie = `access_token=logout`;
    location.reload();
  };

  return (
    <Background className={styles.Background}>
      <Logout
        onClick={clearCookie}
        style={{
          position: "absolute",
          right: 10,
          bottom: 10,
          cursor: "pointer",
        }}
        size={40}
        color="var(--color-neutrals-n-600)"
      />
      {transition((style, state) => {
        return userInf ? (
          <a.div style={style} className={styles.HomeContainer}>
            {state === HomeTypes.Overview ? (
              <HomeOverview userInfo={userInf} setCurrState={setCurrState} />
            ) : state === HomeTypes.EditProfile ? (
              <EditProfile userInf={userInf} />
            ) : (
              <CreateProject setCurrState={setCurrState} />
            )}
          </a.div>
        ) : (
          <SpinningLoading size={60} />
        );
      })}
    </Background>
  );
};

type OverViewProps = {
  userInfo: User;
  setCurrState: (val: HomeTypes) => void;
};

const HomeOverview: React.FC<OverViewProps> = (props) => {
  const breakPoint = useBreakPoints();
  const projects = [];
  console.log(props.userInfo);

  return (
    <div className={styles.HomeContentWrapper}>
      <div className={styles.InfoContainer}>
        <img
          src={props.userInfo.profilePic}
          className={classNames(styles.img)}
        />
        <div className={styles.UserInitials}>
          <Text1 text={props.userInfo.username} />
          <Text3 className={styles.Email} text={props.userInfo.email} />
        </div>
        <ReqButton1
          className={styles.EditProfileBtn}
          text="Edit Profile"
          onClick={() => props.setCurrState(HomeTypes.EditProfile)}
        />
      </div>
      <div className={styles.ProjectSection}>
        <Text2
          text="Projects"
          style={{ fontWeight: "bold", fontSize: "xx-large" }}
        />
        <div className={styles.ProjectsList}>
          <div
            className={classNames(styles.ProjCard, styles.AddProj, {
              [styles.AddMax]: projects.length < 4,
            })}
            onClick={() => props.setCurrState(HomeTypes.CreateProject)}
          >
            <Plus
              size={breakPoint > Breakpoints.Small ? 30 : 20}
              color="var(--color-button)"
            />
          </div>
          <ProjCard project={sProj1} addMax />

          <ProjCard project={sProj1} addMax />
          <ProjCard project={sProj1} addMax />
          <ProjCard project={sProj1} addMax />
        </div>
      </div>
    </div>
  );
};

type EditProfileProps = {
  userInf: User;
};

const EditProfile: React.FC<EditProfileProps> = (props) => {
  return (
    <div className={styles.EditProfileContentWrapper}>
      <img src={props.userInf.profilePic} className={classNames(styles.img)} />
      <ReqButton1 text="Save" />
    </div>
  );
};

type CreateProjInfo = {
  title: string;
  picture: string;
  description: string;
};

type CreajeProjectProps = {
  setCurrState: (val: HomeTypes) => void;
};

const CreateProject: React.FC<CreajeProjectProps> = (props) => {
  const [info, setInfo] = useState<CreateProjInfo>({
    title: "",
    picture: "/public//defaultProjPFP.png",
    description: ""
  });
  const { isPending } = useRequestStates(RequestTypes.CreateProject);
  const ref = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | undefined>();

  const onFileClick = () => {
    ref.current && ref.current.click();
  };

  const onChange = (event: React.ChangeEvent) => {
    const target = event.target as HTMLInputElement;
    const selectedFile = target.files && target.files[0];
    const reader = new FileReader();
    reader.onload = function (event) {
      setInfo((prev) => ({ ...prev, picture: event.target?.result?.toString() ?? prev.picture }));
    };
    selectedFile && reader.readAsDataURL(selectedFile);
  };

  const changeUserName = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    setInfo((prev) => ({...prev, title: event.target.value}))
  }, [])

  const changeDescription = useCallback((event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInfo((prev) => ({...prev, description: event.target.value}))
  }, [])
  const dispatch = useDispatch();
  const handleConfirm = useCallback(() => {
    dispatch(ReqActions.setState({
      requestState: RequestState.Pending,
      reqType: RequestTypes.CreateProject
    }))
    dispatch(ProjectActions.createProject(info))
  }, [])

  return (
    <div className={styles.CreateProjectContentWrapper}>
      <Text1 text="Create new project" />
      <div className={styles.PicWrapper}>
        <img src={info.picture} className={classNames(styles.CreateProjImg)} />
        <div onClick={onFileClick} className={classNames(styles.ImgCover)}>
          Add Picture
        </div>
        <input
          accept=".jpg, .png, .jpeg"
          onChange={onChange}
          ref={ref}
          type="file"
          name="fileInp"
          id=""
          hidden
        />
      </div>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: "20px",
          width: "100%",
        }}
      >
        <UsernameInputBar
          placeholder="Project Name"
          value={info.title}
          onChange={changeUserName}
        />
        <TextAreaInputBar
        onChange={changeDescription
        }
          value={info.description}
          style={{
            maxHeight: "300px",
            maxWidth: "100%",
            minWidth: "100%",
            minHeight: "80px",
          }}
          placeholder="Description"
        />
      </div>
      <div className={styles.ButtonsWrapper}>
        <HollowButton
          text="Back"
          onClick={() => props.setCurrState(HomeTypes.Overview)}
          style={{ width: 100 }}
        />
        <ReqButton1 text="Create Project" style={{ width: 150 }} isPending={isPending} onClick={handleConfirm}/>
      </div>
    </div>
  );
};

type ProjCardProps = {
  project: SmallProject;
  addMax?: boolean;
};

const ProjCard: React.FC<ProjCardProps> = (props) => {
  const navigate = useNavigate();
  return (
    <div
      className={classNames(styles.ProjCard, { [styles.AddMax]: props.addMax })}
      onClick={() => navigate("/project/" + props.project.project_id)}
    >
      <img
        className={styles.ProjCardImg}
        loading="lazy"
        src={props.project.picture}
        alt=""
      />
      <div className={styles.ProjInf}>
        <Text2 className={styles.ProjCardName} text={props.project.title} />
        <div className={styles.MembersInf}>
          <Text3 text={props.project.members_number} />
          <Members size={15} color="var(--color-neutrals-n-600)" />
        </div>
      </div>
    </div>
  );
};
