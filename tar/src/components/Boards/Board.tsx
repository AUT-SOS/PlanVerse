import { useNavigate, useParams } from "react-router-dom";
import styles from "./Board.module.scss";
import { project1 } from "../../utils/testCase";
import { Title } from "../../ui/Title";
import React, {
  PropsWithChildren,
  useCallback,
  useEffect,
  useRef,
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
import { Member, RequestState, RequestTypes, Task } from "../../utils/types";
import { AssignMemberItem, MemberItem } from "../../ui/MemberItem";
import {
  HollowButton,
  ReqButton,
  ReqButton1,
  ReqButtonWithIcon,
} from "../../ui/ReqButton";
import { AddMember } from "../../ui/Icons/AddMember";
import { Modal } from "../../ui/Modal";
import { ReqActions } from "../../redux/slices/req.slice";
import {
  UsernameInputBar,
  TextAreaInputBar,
  EmailInputBar,
  InputBar,
} from "../../ui/InputBar";
import { useRequestStates } from "../../utils/hooks";
import { validateEmail } from "../../utils/regex";
import { showFailToastMessage } from "../../main";
import { CreateProject } from "../Home/Home";
import { TaskBoard } from "./TaskBoard";
import { Logout } from "../../ui/Icons/Logout";

type Props = {};

enum SliderTypes {
  Members,
}

enum ModalType {
  AddMember,
  ProjectSetting,
  TaskInfo,
}

export const Board: React.FC<Props> = (props) => {
  const projId = useParams().id;
  const navigate = useNavigate();
  const [showModal, setShowModal] = useState<ModalType | undefined>(undefined);

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

  const transitionModal = useTransition(showModal, {
    from: {
      opacity: "0",
    },
    enter: {
      opacity: "1",
    },
    leave: {
      opacity: "0",
    },
  });

  const dispatch = useDispatch();
  const { project, members, myId, states } = useSelector(
    (state: RootState) => ({
      project: state.project.fullProject,
      members: state.project.members,
      myId: state.auth.myId,
      states: state.project.states,
    })
  );

  const amIAdmin = members?.find((item) => item.id == myId)?.is_admin;

  const [taskInf, setTaskInf] = useState<
    { task: Task; state_id: string } | undefined
  >(undefined);

  useEffect(() => {
    dispatch(ProjectActions.getFullProject(projId));
  }, []);
  const [visible, setVisible] = useState(false);
  const [sliderContent, setSliderContent] = useState(SliderTypes.Members);

  const handleSettings = useCallback(() => {
    setShowModal(ModalType.ProjectSetting);
  }, []);

  const handleMembers = useCallback(() => {
    dispatch(ProjectActions.getFullProject(projId));
    setSliderContent(SliderTypes.Members);
    setVisible(true);
  }, []);

  const handleOpenTask = (task: Task, state_id: string) => {
    setShowModal(ModalType.TaskInfo);
    setTaskInf({ task, state_id });
  };

  const sliderTitle =
    sliderContent === SliderTypes.Members ? "Members" : "Task Details";

  const amIOwner = myId == project?.owner_id;

  return project && members ? (
    <div className={styles.BoardWrapper}>
      <MoreInfoSlider
        title={sliderTitle}
        visible={visible}
        setVisible={setVisible}
        showInviteIcon={sliderTitle === "Members"}
        openModal={setShowModal}
      >
        {sliderContent === SliderTypes.Members ? (
          <ProjectMembers
            members={members}
            ownerId={project.owner_id}
            projectId={project.project_id}
          />
        ) : (
          <></>
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
            {amIAdmin && (
              <Settings
                onClick={handleSettings}
                size={30}
                color={"white"}
                className={styles.Icon}
              />
            )}
            {!amIOwner && (
              <Logout
                size={30}
                color={"white"}
                className={styles.Icon}
                onClick={() => {
                  dispatch(
                    ProjectActions.deleteProject({
                      project_id: projId,
                      isDelete: false,
                    })
                  );
                }}
              />
            )}
            <Home
              size={30}
              color={"white"}
              className={styles.Icon}
              onClick={() => navigate("/home")}
            />
          </div>
        </a.div>
      ))}

      <div className={styles.BoardContent}>
        {states && (
          <TaskBoard
            openTask={handleOpenTask}
            projectId={projId}
            states={states}
          />
        )}
      </div>

      {transitionModal((style, state) =>
        state != undefined ? (
          <a.div
            style={style}
            className={classNames(styles.Modal, { [styles.show]: showModal })}
            onClick={() => setShowModal(undefined)}
          >
            {state === ModalType.AddMember ? (
              <AddMemberFC
                projId={project.project_id}
                closeModal={setShowModal}
                members={members}
              />
            ) : state === ModalType.ProjectSetting ? (
              <CreateProject
                handleClose={() => setShowModal(undefined)}
                projectInf={{ ...project }}
                projId={projId}
                className={styles.EditProjWrapper}
                amIOwner={project.owner_id == myId}
              />
            ) : state === ModalType.TaskInfo ? (
              <TaskInfo
                closeModal={() => setShowModal(undefined)}
                projId={projId}
                className={styles.EditProjWrapper}
                task={taskInf!.task}
                state_id={taskInf!.state_id}
                members={members}
                amIAdmin={Boolean(amIAdmin)}
              />
            ) : (
              <></>
            )}
          </a.div>
        ) : (
          <></>
        )
      )}
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
  showInviteIcon?: boolean;
  openModal: (val?: ModalType) => void;
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
        {props.showInviteIcon && (
          <ReqButtonWithIcon
            onClick={() => props.openModal(ModalType.AddMember)}
            text="Add Member"
            style={{ width: "auto", minWidth: 0 }}
            title="Add Member"
          >
            <AddMember color="white" size={20} />
          </ReqButtonWithIcon>
        )}
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
  const { myId } = useSelector((state: RootState) => ({
    myId: state.auth.myId,
  }));
  const amIAdmin = props.members.find((item) => item.id == myId)?.is_admin;

  const sortedMembers = [...props.members];
  sortedMembers.sort((item) => {
    return item.is_admin ? -1 : 1;
  });

  return (
    <div className={styles.MemberList}>
      {sortedMembers.map((item, index) => {
        return (
          <MemberItem
            amIAdmin={amIAdmin}
            ownerId={props.ownerId}
            amIOwner={props.ownerId == myId}
            member={item}
            key={index}
            projectId={props.projectId}
            myId={myId}
          />
        );
      })}
    </div>
  );
};

type AddMemberProps = {
  closeModal: (val?: ModalType) => void;
  projId: string;
  members: Member[];
};

const AddMemberFC: React.FC<AddMemberProps> = (props) => {
  const [emails, setEmails] = useState<Set<string>>(new Set());
  const { isPending } = useRequestStates(RequestTypes.CreateProject);
  const [emailInp, setEmailInp] = useState<string>("");

  const handleEmailChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setEmailInp(event.target.value);
    },
    []
  );

  const handleAddEmail = useCallback(
    (event: React.KeyboardEvent) => {
      if (validateEmail(emailInp) && event.key == "Enter") {
        if (props.members.find((item) => emailInp == item.email)) {
          showFailToastMessage("User with this email has joined");
          return;
        }
        const tempSet = emails;
        tempSet.add(emailInp);
        setEmails(tempSet);
        setEmailInp("");
      }
    },
    [emailInp, props.members, emails]
  );

  const handleDeleteEmail = useCallback(
    (email: string) => {
      const tempSet = emails;
      tempSet.delete(email);
      setEmails(tempSet);

      setEmails(emails);
    },
    [emails]
  );

  const dispatch = useDispatch();
  const handleConfirm = useCallback(() => {
    dispatch(
      ReqActions.setState({
        requestState: RequestState.Pending,
        reqType: RequestTypes.ShareLink,
      })
    );
    dispatch(
      ProjectActions.shareLink({ id: props.projId, emails: Array.from(emails) })
    );
    setEmails(new Set());
  }, [emails, props.projId]);

  return (
    <div
      className={styles.AddMemberContentWrapper}
      onClick={(e) => e.stopPropagation()}
    >
      <Text1 text="Invite new members" />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: "20px",
          width: "100%",
        }}
      >
        <EmailInputBar
          placeholder="Enter a vaild email to invite (ENTER to submit)"
          value={emailInp}
          onChange={handleEmailChange}
          onKeyDown={handleAddEmail}
        />
        <div className={styles.EmailsWrapper}>
          {Array.from(emails).map((item, index) => (
            <div key={index} className={styles.EmailItem}>
              <>
                <div
                  className={styles.EmailDeleteHandler}
                  onClick={() => handleDeleteEmail(item)}
                >
                  X
                </div>
                {item}
              </>
            </div>
          ))}
        </div>
      </div>
      <div className={styles.ButtonsWrapper}>
        <HollowButton
          text="Back"
          onClick={() => props.closeModal(undefined)}
          style={{ width: 100 }}
        />
        <ReqButton1
          text="Send Invite Link"
          style={{ width: 150 }}
          isPending={isPending}
          onClick={handleConfirm}
          disable={emails.size < 1}
        />
      </div>
    </div>
  );
};

type TaskInfoProps = React.HTMLProps<HTMLDivElement> & {
  task: Task;
  projId: string;
  closeModal: (val?: ModalType) => void;
  state_id: string;
  members: Member[];
  amIAdmin: boolean;
};

const TaskInfo: React.FC<TaskInfoProps> = (props) => {
  const [task, setTask] = useState<Task>(props.task);

  const { isPending } = useRequestStates(RequestTypes.CreateProject);

  const handleTitleChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setTask((prev) => ({ ...prev, title: event.target.value }));
    },
    []
  );

  const [assigned, setAssign] = useState<number[]>(props.task.performers ?? []);

  const handleClose = () => {
    dispatch(ProjectActions.setTask(undefined));
    props.closeModal(undefined);
  };

  const handleDescriptionChange = useCallback(
    (event: React.ChangeEvent<HTMLTextAreaElement>) => {
      setTask((prev) => ({ ...prev, description: event.target.value }));
    },
    []
  );

  const handleMemberClick = (member: Member, isAssign: boolean) => {
    if (!props.amIAdmin) return;
    if (isAssign) {
      setAssign((prev) => [...prev, Number(member.id)]);
      dispatch(
        ProjectActions.assign({
          project_id: props.projId,
          task_id: props.task.task_id,
          performer_id: member.id,
          isAdd: true,
        })
      );
    } else {
      setAssign((prev) => prev.filter((item) => item != Number(member.id)));
      dispatch(
        ProjectActions.assign({
          project_id: props.projId,
          task_id: props.task.task_id,
          performer_id: member.id,
          isAdd: false,
        })
      );
    }
  };

  const dispatch = useDispatch();
  const handleConfirm = () => {
    dispatch(
      ReqActions.setState({
        requestState: RequestState.Pending,
        reqType: RequestTypes.EditTask,
      })
    );
    dispatch(
      ProjectActions.editTask({
        ...task,
        project_id: props.projId,
        state_id: props.state_id,
      })
    );
  };

  const handleDeleteTask = () => {
    dispatch(
      ProjectActions.deleteTask({
        project_id: props.projId,
        task_id: props.task.task_id,
      })
    );
    handleClose();
  };
  const filteredMembers = [...props.members];

  return task ? (
    <div
      className={styles.AddMemberContentWrapper}
      onClick={(e) => e.stopPropagation()}
    >
      <Text1 text="Task Info" />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: "20px",
          width: "100%",
        }}
      >
        <InputBar
          placeholder="Task Title"
          value={task.title}
          onChange={handleTitleChange}
        />
        <TextAreaInputBar
          onChange={handleDescriptionChange}
          value={task.description}
          style={{
            maxHeight: "300px",
            maxWidth: "100%",
            minWidth: "100%",
            minHeight: "80px",
          }}
          placeholder="Description"
        />
        <div className={styles.AssignList}>
          {filteredMembers.map((item) => {
            const isAssigned =
              assigned.find((a) => a == Number(item.id)) !== undefined;

            return props.amIAdmin ? (
              <AssignMemberItem
                key={item.id}
                isAssigned={
                  assigned.find((a) => a == Number(item.id)) !== undefined
                }
                member={item}
                setAssign={handleMemberClick}
              />
            ) : (
              isAssigned && (
                <AssignMemberItem
                  key={item.id}
                  isAssigned={
                    assigned.find((a) => a == Number(item.id)) !== undefined
                  }
                  member={item}
                  setAssign={handleMemberClick}
                />
              )
            );
          })}
        </div>
      </div>
      <div className={styles.ButtonsWrapper}>
        <div style={{ display: "flex", gap: 10 }}>
          <HollowButton
            text="Back"
            onClick={handleClose}
            style={{ width: 100 }}
          />
          {props.amIAdmin && (
            <ReqButton
              style={{ borderRadius: 10 }}
              text="Delete Task"
              onClick={handleDeleteTask}
            />
          )}
        </div>
        {props.amIAdmin && (
          <ReqButton1
            text="Submit Changes"
            style={{ width: 150 }}
            isPending={isPending}
            onClick={handleConfirm}
            disable={!props.amIAdmin}
          />
        )}
      </div>
    </div>
  ) : (
    <SpinningLoading size={50} />
  );
};
