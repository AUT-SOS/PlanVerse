import classNames from "classnames";
import React from "react";
import styles from "./InputBar.module.scss";

type Props = React.HTMLProps<HTMLInputElement>;

export const InputBar: React.FC<Props> = (props) => {
  return (
    <input
      type="text"
      {...props}
      className={classNames(props.className, styles.commonInput)}
    />
  );
};

export const PasswordInputBar: React.FC<Props> = (props) => {
  return (
    <input
      type="password"
      {...props}
      className={classNames(
        props.className,
        styles.commonInput,
        styles.commonInput,
        {
          [styles.error]: !validatePassword(props.value as string),
        }
      )}
    />
  );
};

export const EmailInputBar: React.FC<Props> = (props) => {
  return (
    <input
      type="text"
      {...props}
      className={classNames(props.className, styles.commonInput, {
        [styles.error]: !validateEmail(props.value as string),
      })}
      value={props.value}
    />
  );
};

const validateEmail = (email?: string) => {
  return email
    ? email.match(
        /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      )
    : true;
};

const validatePassword = (password?: string) => {
  return password
    ? password.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/)
    : true;
};
