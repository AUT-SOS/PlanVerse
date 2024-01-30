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
import { Member, RequestState, RequestTypes } from "../../utils/types";
import { MemberItem } from "../../ui/MemberItem";
import {
  HollowButton,
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
} from "../../ui/InputBar";
import { useRequestStates } from "../../utils/hooks";
import { validateEmail } from "../../utils/regex";
import { showFailToastMessage } from "../../main";
import { CreateProject } from "../Home/Home";

type Props = {};

enum SliderTypes {
  Settings,
  Members,
}

enum ModalType {
  AddMember,
  ProjectSetting,
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
  const { project, members, myId } = useSelector((state: RootState) => ({
    project: state.project.fullProject,
    members: state.project.members,
    myId: state.auth.myId
  }));

  const amIAdmin = members?.find((item) => item.id == myId)?.is_admin;

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

  const sliderTitle =
    sliderContent === SliderTypes.Members ? "Members" : "Project Setting";

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
            {amIAdmin && <Settings
              onClick={handleSettings}
              size={30}
              color={"white"}
              className={styles.Icon}
            />}
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
  console.log(amIAdmin, myId, props.ownerId);

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
