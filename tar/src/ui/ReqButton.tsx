import React from "react";
import styles from "./ReqButton.module.scss";
import classNames from "classnames";
import { SpinningLoading } from "./SpinningLoading";

type Props = React.HTMLProps<HTMLDivElement> & {
  isPending?: boolean;
  text: string;
};

export const ReqButton: React.FC<Props> = (props) => {
  return (
    <div
      {...props}
      onClick={props.onClick}
      className={classNames(
        styles.ButtonWrapper,
        { [styles.isPending]: props.isPending },
        props.className
      )}
    >
      {props.isPending ? <SpinningLoading size={20} /> : props.text}
    </div>
  );
};
