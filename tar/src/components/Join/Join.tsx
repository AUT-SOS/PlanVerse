import React from "react";
import { useParams } from "react-router-dom";
import { ReqButton } from "../../ui/ReqButton";
import classNames from "classnames";
import { Members } from "../../ui/Icons/Members";
import "./Join.scss";
import { Project } from "../../utils/types";
import { project1 } from "../../utils/testCase";
import { HollowButton } from "../../ui/HollowButton";
import { Title } from "../../ui/Title";

export const Join: React.FC = (props) => {
  const params = useParams();
  const data = project1;

  return (
    <>
      <Title text="PlanVerse" href="/"/>
      <div className={classNames("JoinWrapper")}>
        <div className={classNames("JoinCard")}>
          <img className="GroupIMG" src={data.background} alt="" />
          <div className="GroupName">{data.name}</div>
          <div className="GroupMemberWrapper">
            <div className="GroupMembers">
              {data.members.slice(0, 3).map((item) => {
                return (
                  <img
                    key={item.id}
                    title={item.username}
                    src={item.profile_pic}
                    className="MemberPrev"
                  />
                );
              })}
            </div>
            <Members size={22} color={"var(--color-neutrals-n-500)"} />
            <p className="MembersCount">{data.members.length}</p>
          </div>
          <div className="ButtonWrapper">
            <HollowButton text="Cancel" width={40} height={20} />
            <ReqButton text="Join Project" width={40} height={20} />
          </div>
        </div>
      </div>
    </>
  );
};
