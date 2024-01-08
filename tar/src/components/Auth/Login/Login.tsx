import React, { useCallback, useState } from "react";
import styles from "./Login.module.scss";
import classNames from "classnames";
import { Text0, Text1, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import {
  EmailInputBar,
  InputBar,
  PasswordInputBar,
} from "../../../ui/InputBar";
import { ReqButton } from "../../../ui/ReqButton";
import { AuthState, LoginForm } from "../../../utils/types";
import { useDispatch, useSelector } from "react-redux";
import { AuthActions } from "../../../redux/slices/auth.slice";
import { RootState } from "../../../redux/store";

export const Login: React.FC = () => {
  const [loginForm, setLoginForm] = useState<LoginForm>({email: "", password: ""});
  const dispatch = useDispatch();

  const authStatus = useSelector((state: RootState) => state.auth.authState);

  const handleLogin = useCallback(() => {
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
        />
        <PasswordInputBar
          placeholder={strings.auth.enterPass}
          value={loginForm.password}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setLoginForm((prev) => ({ ...prev, password: e.target.value }))
          }
        />
      </div>
      <ReqButton
        isPending={authStatus === AuthState.Pending}
        text={strings.auth.login}
        style={{ fontSize: "medium", width: "40%" }}
        onClick={handleLogin}
      />
      <Text3
        text={strings.auth.forgotPass}
        style={{ color: "var(--color-button)", cursor: "pointer" }}
      />
    </div>
  );
};
