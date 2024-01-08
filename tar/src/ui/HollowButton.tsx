import React from "react";
import styles from "./HollowButton.module.scss";
import classNames from "classnames";


type Props = React.HTMLProps<HTMLDivElement> & {
  isPending?: boolean;
  text: string;
  width: number;
  height: number;
};

export const HollowButton : React.FC<Props> = (props) => {
  return (
    <div
      onClick={props.onClick}
      style={{
        paddingTop: props.height / 2,
        paddingBottom: props.height / 2,
        paddingLeft: props.width / 2,
        paddingRight: props.width / 2,
      }}
      className={classNames(styles.ButtonWrapper)}
    >
      {props.text}
    </div>
  );
};
