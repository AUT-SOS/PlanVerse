import classNames from "classnames";
import React, { forwardRef } from "react";
import styles from "./InputBar.module.scss";
import {
  validatePassword,
  validateEmail,
  validateUsername,
} from "../utils/regex";
import strings from "../utils/text";

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
      title={strings.auth.password_lim}
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

export const UsernameInputBar: React.FC<Props> = (props) => {
  return (
    <input
      type="text"
      title={strings.auth.username_lim}
      {...props}
      className={classNames(props.className, styles.commonInput, {
        [styles.error]: !validateUsername(props.value as string),
      })}
      value={props.value}
    />
  );
};

export const DigitInput = forwardRef<
  HTMLInputElement,
  Props & {
    onChangeInput: (e: React.ChangeEvent<HTMLInputElement>) => void;
  }
>((props, ref) => {
  return (
    <input
      ref={ref}
      onChange={props.onChangeInput}
      type="number"
      className={classNames(props.className, styles.digitInput)}
      {...props}
    />
  );
});
