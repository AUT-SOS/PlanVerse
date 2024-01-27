import classNames from "classnames";
import React, { forwardRef } from "react";
import styles from "./InputBar.module.scss";
import {
  validatePassword,
  validateEmail,
  validateUsername,
} from "../utils/regex";
import strings from "../utils/text";

type Props = React.HTMLProps<HTMLInputElement>  & {
  error?: boolean
};

export const InputBar: React.FC<Props> = (props) => {
  return (
    <input
      type="text"
      {...props}
      className={classNames(props.className, styles.commonInput, {
        [styles.error]: props.error,
      })}
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
          [styles.error]: props.error || !validatePassword(props.value as string),
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
        [styles.error]: props.error || !validateEmail(props.value as string),
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
        [styles.error]: props.error || !validateUsername(props.value as string),
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
      className={classNames(props.className, styles.digitInput, {
        [styles.error]: props.error
      })}
      {...props}
    />
  );
});

export const TextAreaInputBar: React.FC<React.HTMLProps<HTMLTextAreaElement>> = (props) => {
  return (
    <textarea
      type="text"
      title={strings.auth.username_lim}
      {...props}
      className={classNames(props.className, styles.commonInput)}
      value={props.value}
    />
  );
};
