import React from "react";
import styles from "./ReqButton.module.scss";
import classNames from "classnames";
import { Spinner } from "./Icons/Spinner";
import { SpinningLoading } from "./SpinningLoading";

type Props = React.HTMLProps<HTMLDivElement> & {
  isPending?: boolean;
  text: string;
  width: number;
  height: number;
};

export const ReqButton: React.FC<Props> = (props) => {
  return (
    <div
      onClick={props.onClick}
      style={{
        paddingTop: props.height / 2,
        paddingBottom: props.height / 2,
        paddingLeft: props.width / 2,
        paddingRight: props.width / 2,
      }}
      className={classNames(styles.ButtonWrapper, {[styles.isPending] : props.isPending}, props.className)}
    >
      {props.isPending ?  <SpinningLoading size={20}/> : props.text}
    </div>
  );
};
