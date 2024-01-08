import { HTMLProps } from "react";
import styles from "./Title.module.scss"
import classNames from "classnames";

type Props = HTMLProps<HTMLParagraphElement> & {
  text: string;
  href?: string
}

export const Title: React.FC<Props> = (props) => {
  return <p {...props} className={classNames(styles.Title, props.className)}>{props.text}</p>;
};
