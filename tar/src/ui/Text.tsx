import classNames from "classnames";
import React from "react";
import styles from "./Text.module.scss";

type Props = React.HTMLProps<HTMLDivElement> & {
      text: string;
}

export const Text0: React.FC<Props> = (props) => {
      return <div {...props} className={classNames(props.className, styles.text0)}>
            {props.text}
      </div>
}

export const Text1: React.FC<Props> = (props) => {
      return <div {...props} className={classNames(props.className, styles.text1)}>
            {props.text}
      </div>
}

export const Text2: React.FC<Props> = (props) => {
      return <div {...props} className={classNames(props.className, styles.text2)}>
            {props.text}
      </div>
}

export const Text3: React.FC<Props> = (props) => {
      return <div {...props} className={classNames(props.className, styles.text3)}>
            {props.text}
      </div>
}

export const Text4: React.FC<Props> = (props) => {
      return <div {...props} className={classNames(props.className, styles.text4)}>
            {props.text}
      </div>
}