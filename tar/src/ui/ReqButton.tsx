import React, { CSSProperties, useEffect } from "react";
import styles from "./ReqButton.module.scss";
import classNames from "classnames";
import { SpinningLoading } from "./SpinningLoading";
import { useShake } from "../utils/hooks";
import { Interpolation, a } from "@react-spring/web";

type Props = Omit<React.HTMLProps<HTMLDivElement>, "ref" | "style"> & {
  isPending?: boolean;
  text: string;
  style?: CSSProperties & { x?: Interpolation<number, number> };
  disable?: boolean;
};

export const ReqButton: React.FC<Props> = (props) => {
  return (
    <a.div
      {...props}
      onClick={props.onClick}
      className={classNames(
        styles.ButtonWrapper,
        { [styles.isPending]: props.isPending, [styles.disable]: props.disable },
        props.className
      )}
      style={{
        ...props.style,
      }}
    >
      {props.isPending ? <SpinningLoading size={20} /> : props.text}
    </a.div>
  );
};

export const ResendEmailButton: React.FC<Props> = (props) => {
  return (
    <a.div
      {...props}
      onClick={props.onClick}
      className={classNames(
        styles.ResendButtonWrapper,
        { [styles.isPending]: props.isPending, [styles.disable]: props.disable },
        props.className
      )}
      style={{
        ...props.style,
      }}
    >
      {props.isPending ? <SpinningLoading size={20} /> : props.text}
    </a.div>
  );
};

export const ReqButton1: React.FC<Props> = (props) => {
  return (
    <a.div
      {...props}
      onClick={props.onClick}
      className={classNames(
        styles.ButtonWrapper1,
        { [styles.isPending]: props.isPending, [styles.disable]: props.disable },
        props.className
      )}
      style={{
        ...props.style,
      }}
    >
      {props.isPending ? <SpinningLoading size={20} /> : props.text}
    </a.div>
  );
};

export const HollowButton: React.FC<Props> = (props) => {
  return (
    <a.div
      {...props}
      onClick={props.onClick}
      className={classNames(
        styles.ButtonWrapperHollow,
        { [styles.isPending]: props.isPending, [styles.disable]: props.disable },
        props.className
      )}
      style={{
        ...props.style,
      }}
    >
      {props.isPending ? <SpinningLoading size={20} /> : props.text}
    </a.div>
  );
};
