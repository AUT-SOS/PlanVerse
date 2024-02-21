import React, { useCallback } from "react";
import { Member } from "../utils/types";
import { Text2, Text3, Text4 } from "./Text";
import styles from "./MemberItem.module.scss";
import classNames from "classnames";
import { Demote } from "./Icons/Demote";
import { Promote } from "./Icons/Promote";
import { shallowEqual, useDispatch } from "react-redux";
import { ProjectActions } from "../redux/slices/project.slice";
import { Star } from "./Icons/Star";
import { a, useTransition } from "@react-spring/web";

type Props = React.HTMLProps<HTMLDivElement> & {
  member: Member;
  ownerId?: string;
  amIAdmin?: boolean;
  amIOwner?: boolean;
  projectId: string;
  myId?: string;
};

export const MemberItemFC: React.FC<Props> = (props) => {
  const pfp =
    props.member.profile_pic.length > 0
      ? props.member.profile_pic
      : "/DefaultPFP.jpg";

  const dispatch = useDispatch();

  const handlePromote = useCallback(() => {
    dispatch(
      ProjectActions.changeMemberRole({
        projectId: props.projectId,
        userId: props.member.id,
        isPromote: true,
      })
    );
  }, [props.projectId, props.member]);

  const handleDemote = useCallback(() => {
    dispatch(
      ProjectActions.changeMemberRole({
        projectId: props.projectId,
        userId: props.member.id,
        isPromote: false,
      })
    );
  }, [props.projectId, props.member]);

  return (
    <div
      {...props}
      className={classNames(styles.MemberItemWrapper, props.className)}
    >
      <img className={styles.MemberImg} src={pfp} alt="" />
      <div className={styles.InfoWrapper}>
        <Text2 text={props.member.username} />
        <Text3 text={props.member.email} />
      </div>

      <div style={{ marginLeft: "auto", display: "flex", gap: 10 }}>
        {props.member.id != props.myId ? (
          props.amIOwner ? (
            props.member.is_admin ? (
              <Demote style={{cursor: "pointer"}} color="red" size={20} onClick={handleDemote} />
            ) : (
              <Promote style={{cursor: "pointer"}} color="green" size={20} onClick={handlePromote} />
            )
          ) : props.amIAdmin && !props.member.is_admin ? (
            <Promote style={{cursor: "pointer"}} color="green" size={20} onClick={handlePromote} />
          ) : (
            <></>
          )
        ) : (
          <></>
        )}
        <Text4
          text={
            props.ownerId === props.member.id
              ? "Owner"
              : props.member.is_admin
              ? "Admin"
              : "Member"
          }
        />
      </div>
    </div>
  );
};

type AssignMemberProps = React.HTMLProps<HTMLDivElement> & {
  member: Member;
  isAssigned: boolean;
  setAssign: (member: Member, isAssign: boolean) => void;
};

export const AssignMemberItem: React.FC<AssignMemberProps> = (props) => {
  const pfp =
    props.member.profile_pic.length > 0
      ? props.member.profile_pic
      : "/DefaultPFP.jpg";

  const handleClick = useCallback(() => {
    props.setAssign(props.member, !props.isAssigned);
  }, [props.setAssign, props.member]);

  const transition = useTransition(props.isAssigned, {
    from:{
      opacity: 0,
      x: 20
    },
    enter:{
      opacity: 1,
      x: 0
    },
    leave: {
      opacity: 0,
      x: 20
    }
  })

  return (
    <div
      {...props}
      className={classNames(styles.MemberItemWrapper, props.className)}
      onClick={handleClick}
      style={{cursor: "pointer"}}
    >
      <img className={styles.MemberImg} src={pfp} alt="" />
      <div className={styles.InfoWrapper}>
        <Text2 text={props.member.username} />
      </div>

      <div style={{ marginLeft: "auto", display: "flex", gap: 10 }}>
        {transition((style, state) => state && <a.div style={style}>
         <Star color="gold" size={25}/>

        </a.div>)}
      </div>
    </div>
  );
};

export const MemberItem = React.memo(MemberItemFC, shallowEqual);
