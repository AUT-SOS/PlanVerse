import React, { useCallback, useEffect, useRef, useState } from "react";
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
  UserEditType,
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
import {
  EmailInputBar,
  PasswordInputBar,
  TextAreaInputBar,
  UsernameInputBar,
} from "../../ui/InputBar";
import { ReqActions } from "../../redux/slices/req.slice";
import { ProjectActions } from "../../redux/slices/project.slice";
import { UserActions } from "../../redux/slices/user.slice";
import {
  validateEmail,
  validatePassword,
  validateUsername,
} from "../../utils/regex";

enum HomeTypes {
  Overview,
  CreateProject,
  EditProfile,
}

export const Home: React.FC = (props) => {
  const { userInf, isAuth, userId } = useSelector((state: RootState) => ({
    userInf: state.users.me,
    isAuth: state.auth.authState === AuthState.Authenticated,
    userId: state.auth.myId,
  }));
  const [currState, setCurrState] = useState<HomeTypes>(HomeTypes.Overview);
  const navigate = useNavigate();
  if (!isAuth && userId === undefined) {
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
              <EditProfile userInf={userInf} setCurrState={setCurrState} />
            ) : (
              <CreateProject
                handleClose={() => setCurrState(HomeTypes.Overview)}
              />
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
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(ProjectActions.getMyProjects());
  }, []);
  const projList =
    useSelector((state: RootState) => state.project.myProjects) ?? [];

  return (
    <div className={styles.HomeContentWrapper}>
      <div className={styles.InfoContainer}>
        <img
          src={props.userInfo.profile_pic}
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
              [styles.AddMax]: projList.length < 4,
            })}
            onClick={() => props.setCurrState(HomeTypes.CreateProject)}
          >
            <Plus
              size={breakPoint > Breakpoints.Small ? 30 : 20}
              color="var(--color-button)"
            />
          </div>
          {projList.map((item, index) => {
            return <ProjCard key={index} project={item} addMax />;
          })}
        </div>
      </div>
    </div>
  );
};

type EditProfileProps = {
  userInf: User;
  setCurrState: (val: HomeTypes) => void;
};

const EditProfile: React.FC<EditProfileProps> = (props) => {
  const [info, setInfo] = useState<UserEditType>({
    ...props.userInf,
    password: "",
  });
  const { isPending } = useRequestStates(RequestTypes.EditUser);
  const ref = useRef<HTMLInputElement>(null);

  const onFileClick = () => {
    ref.current && ref.current.click();
  };

  const changeImg = (event: React.ChangeEvent) => {
    const target = event.target as HTMLInputElement;
    const selectedFile = target.files && target.files[0];
    const reader = new FileReader();
    reader.onload = function (event) {
      setInfo((prev) => ({
        ...prev,
        profile_pic: event.target?.result?.toString() ?? prev.profile_pic,
      }));
    };
    selectedFile && reader.readAsDataURL(selectedFile);
  };

  const changeUsername = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInfo((prev) => ({ ...prev, username: event.target.value }));
    },
    []
  );

  const changeEmail = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInfo((prev) => ({ ...prev, email: event.target.value }));
    },
    []
  );

  const changePassword = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInfo((prev) => ({ ...prev, password: event.target.value }));
    },
    []
  );

  const dispatch = useDispatch();
  const handleConfirm = useCallback(() => {
    dispatch(
      ReqActions.setState({
        requestState: RequestState.Pending,
        reqType: RequestTypes.EditUser,
      })
    );
    dispatch(UserActions.editUserInfo(info));
  }, [info]);

  const submitDisable =
    info.password.length < 8 ||
    (info.email === props.userInf.email &&
      info.profile_pic === props.userInf.profile_pic &&
      info.username === props.userInf.username) ||
    !validateEmail(info.email) ||
    !validateUsername(info.username) ||
    !validatePassword(info.password);

  return (
    <div className={styles.CreateProjectContentWrapper}>
      <Text1 text="Edit Profile" />
      <div className={styles.EditUserPicWrapper}>
        <img
          src={info.profile_pic}
          className={classNames(styles.EditUserProfilePic)}
        />
        <div onClick={onFileClick} className={classNames(styles.ImgCover)}>
          Add Picture
        </div>
        <input
          accept=".jpg, .png, .jpeg"
          onChange={changeImg}
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
          placeholder="Username"
          value={info.username}
          onChange={changeUsername}
        />
        <EmailInputBar
          onChange={changeEmail}
          value={info.email}
          placeholder="Email"
        />
        <PasswordInputBar
          value={info.password}
          placeholder="Enter your password"
          onChange={changePassword}
        />
      </div>
      <div className={styles.ButtonsWrapper}>
        <HollowButton
          text="Back"
          onClick={() => props.setCurrState(HomeTypes.Overview)}
          style={{ width: 100 }}
        />
        <ReqButton1
          text="Submit"
          style={{ width: 150 }}
          isPending={isPending}
          onClick={handleConfirm}
          disable={submitDisable}
        />
      </div>
    </div>
  );
};

export type CreateProjInfo = {
  title: string;
  picture: string;
  description: string;
};

type CreajeProjectProps = {
  handleClose: () => void;
  projectInf?: CreateProjInfo;
  projId?: string;
  className?: string;
};

export const CreateProject: React.FC<CreajeProjectProps> = (props) => {
  const [info, setInfo] = useState<CreateProjInfo>({
    title: props.projectInf?.title ?? "",
    picture: props.projectInf?.picture ?? "/public//defaultProjPFP.png",
    description: props.projectInf?.description ?? "",
  });
  const { isPending } = useRequestStates(RequestTypes.CreateProject);
  const ref = useRef<HTMLInputElement>(null);

  const onFileClick = () => {
    ref.current && ref.current.click();
  };

  const onChange = (event: React.ChangeEvent) => {
    const target = event.target as HTMLInputElement;
    const selectedFile = target.files && target.files[0];
    const reader = new FileReader();
    reader.onload = function (event) {
      setInfo((prev) => ({
        ...prev,
        picture: event.target?.result?.toString() ?? prev.picture,
      }));
    };
    selectedFile && reader.readAsDataURL(selectedFile);
  };

  const changeUserName = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInfo((prev) => ({ ...prev, title: event.target.value }));
    },
    []
  );

  const changeDescription = useCallback(
    (event: React.ChangeEvent<HTMLTextAreaElement>) => {
      setInfo((prev) => ({ ...prev, description: event.target.value }));
    },
    []
  );
  const dispatch = useDispatch();
  const handleConfirm = useCallback(() => {
    if (props.projId) {
      dispatch(
        ReqActions.setState({
          requestState: RequestState.Pending,
          reqType: RequestTypes.EditProject,
        })
      );
      dispatch(ProjectActions.editProject({...info, id: props.projId}));
    } else {
      dispatch(
        ReqActions.setState({
          requestState: RequestState.Pending,
          reqType: RequestTypes.CreateProject,
        })
      );
      dispatch(ProjectActions.createProject(info));
    }
  }, [info, props.projId]);

  return (
    <div
      className={classNames(styles.CreateProjectContentWrapper, props.className)}
      onClick={(e) => e.stopPropagation()}
    >
      <Text1 text={props.projectInf ? "Edit Project" : "Create new project"} />
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
          onChange={changeDescription}
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
          onClick={props.handleClose}
          style={{ width: 100 }}
        />
        <ReqButton1
          text={props.projectInf ? "Edit Project" : "Create Project"}
          style={{ width: 150 }}
          isPending={isPending}
          onClick={handleConfirm}
          disable={info.title.length < 3}
        />
      </div>
    </div>
  );
};

type ProjCardProps = {
  project: SmallProject;
  addMax?: boolean;
  key: number;
};

const ProjCard: React.FC<ProjCardProps> = (props) => {
  const navigate = useNavigate();
  return (
    <div
      key={props.key}
      className={classNames(styles.ProjCard, { [styles.AddMax]: props.addMax })}
      onClick={() => navigate("/projects/" + props.project.project_id)}
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
