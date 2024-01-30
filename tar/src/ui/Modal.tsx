import classNames from "classnames";
import React, { PropsWithChildren } from "react";
import styles from "./Modal.module.scss"

export const Modal: React.FC<PropsWithChildren & React.HTMLProps<HTMLDialogElement>> = (props) => {
      return <dialog {...props} className={classNames(styles.Modal, props.className)}>
            {props.children}
      </dialog>
}