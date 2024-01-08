import React, { PropsWithChildren } from "react";
import styles from "./Background.module.scss"
import classNames from "classnames";

type Props = PropsWithChildren & React.HTMLAttributes<HTMLDivElement>


export const Background : React.FC<Props> = (props) => {
      return <div className={classNames(props.className, styles.bgWrapper)}>
            {props.children}
      </div>
}