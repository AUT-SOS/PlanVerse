import React from "react";
import { Spinner } from "./Icons/Spinner";
import styles from "./SpinningLoading.module.scss";

type Props = {
  size: number;
};

export const SpinningLoading: React.FC<Props> = (props) => {
  return <Spinner size={props.size} className={styles.Spinner}/>;
};
