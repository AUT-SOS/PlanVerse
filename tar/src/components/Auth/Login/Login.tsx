import React, { useCallback, useState } from "react";
import styles from "./Login.module.scss";
import classNames from "classnames";
import { Text1, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import {
  EmailInputBar,
  PasswordInputBar,
} from "../../../ui/InputBar";
import { ReqButton } from "../../../ui/ReqButton";
import {
  LoginForm,
  RequestState,
  RequestTypes,
} from "../../../utils/types";
import { useDispatch } from "react-redux";
import { AuthActions } from "../../../redux/slices/auth.slice";
import { useRequestStates, useShake } from "../../../utils/hooks";
import { a } from "@react-spring/web";
import { validateEmail, validatePassword } from "../../../utils/regex";
import { ReqActions } from "../../../redux/slices/req.slice";

export const Login: React.FC = () => {
  const [loginForm, setLoginForm] = useState<LoginForm>({
    email: "",
    password: "",
  });
  const dispatch = useDispatch();
  const shakeAnimation = useShake(0, 2);

  const { isPending, errorState } = useRequestStates(RequestTypes.Login);


  const handleLogin = useCallback(() => {
    if (
      !loginForm.email ||
      !loginForm.password ||
      !validateEmail(loginForm.email) ||
      !validatePassword(loginForm.password)
    ) {
      shakeAnimation.api.start({
        from: {
          x: 1.5,
        },
        to: {
          x: 0,
        },
      });
      return;
    }
    dispatch(
      ReqActions.setState({
        requestState: RequestState.Pending,
        reqType: RequestTypes.Login,
      })
    );
    dispatch(AuthActions.login(loginForm));
  }, [loginForm]);
  return (
    <div className={classNames(styles.loginWrapper)}>
      <Text1
        text={strings.auth.login}
        style={{ color: "var(--color-neutrals-n-500)" }}
      />
      <div className={classNames(styles.inputWrapper)}>
        <EmailInputBar
          placeholder={strings.auth.enterEmail}
          value={loginForm.email}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setLoginForm((prev) => ({ ...prev, email: e.target.value }))
          }
          error={errorState}
        />
        <PasswordInputBar
          placeholder={strings.auth.enterPass}
          value={loginForm.password}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setLoginForm((prev) => ({ ...prev, password: e.target.value }))
          }
          error={errorState}
        />
      </div>
      <a.div
        style={{
          x: shakeAnimation.x,
          width: "100%",
          display: "flex",
          justifyContent: "center",
        }}
      >
        <ReqButton
          isPending={isPending}
          text={strings.auth.login}
          style={{ fontSize: "medium", width: "40%" }}
          onClick={handleLogin}
        />
      </a.div>
      <Text3
        text={strings.auth.forgotPass}
        style={{ color: "var(--color-button)", cursor: "pointer" }}
      />
    </div>
  );
};
