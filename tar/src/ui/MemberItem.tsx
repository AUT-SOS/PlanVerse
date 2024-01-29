import React, { useCallback } from "react";
import { Member } from "../utils/types";
import { Text2, Text3, Text4 } from "./Text";
import styles from "./MemberItem.module.scss";
import classNames from "classnames";
import { Demote } from "./Icons/Demote";
import { Promote } from "./Icons/Promote";
import { shallowEqual, useDispatch } from "react-redux";
import { ProjectActions } from "../redux/slices/project.slice";

type Props = React.HTMLProps<HTMLDivElement> & {
  member: Member;
  ownerId?: string;
  amIAdmin?: boolean;
  amIOwner?: boolean;
  projectId: string;
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

      <div style={{ marginLeft: "auto", display: "flex", gap: "5" }}>
        {props.amIOwner ? (
          props.member.is_admin ? (
            <Demote size={20} onClick={handleDemote}/>
          ) : (
            <Promote size={20} onClick={handlePromote}/>
          )
        ) : props.amIAdmin && !props.member.is_admin ? (
          <Promote size={20} onClick={handlePromote}/>
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

export const MemberItem = React.memo(MemberItemFC, shallowEqual);