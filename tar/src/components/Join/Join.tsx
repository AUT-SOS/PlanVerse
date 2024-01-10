import React from "react";
import { useParams } from "react-router-dom";
import { ReqButton } from "../../ui/ReqButton";
import classNames from "classnames";
import { Members } from "../../ui/Icons/Members";
import "./Join.module.scss";
import { Project } from "../../utils/types";
import { project1 } from "../../utils/testCase";
import { HollowButton } from "../../ui/HollowButton";
import { Title } from "../../ui/Title";
import { Background } from "../../ui/BackGround";
import styles from "./Join.module.scss";
import strings from "../../utils/text";

export const Join: React.FC = (props) => {
  const params = useParams();
  const data = project1;

  return (
    <>
      <Title text={strings.palverse} href="/" />
      <Background className={styles.JoinWrapper}>
        <div className={classNames(styles.JoinCard)}>
          <div className={styles.contentWrapper}>
            <img className={styles.GroupIMG} src={data.background} alt="" />
            <div className={styles.GroupName}>{data.name}</div>
            <div className={styles.GroupMemberWrapper}>
              <div className={styles.GroupMembers}>
                {data.members.slice(0, 3).map((item) => {
                  return (
                    <img
                      key={item.id}
                      title={item.username}
                      src={item.profile_pic}
                      className={styles.MemberPrev}
                    />
                  );
                })}
              </div>
              <Members size={22} color={"var(--color-neutrals-n-500)"} />
              <p className={styles.MembersCount}>{data.members.length}</p>
            </div>
            <ReqButton
              text={strings.join.joinProject}
              style={{ width: "100%" }}
            />
          </div>
        </div>
      </Background>
    </>
  );
};
