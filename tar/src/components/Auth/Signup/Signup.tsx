import React, { useCallback, useState } from "react";
import styles from "./Signup.module.scss";
import classNames from "classnames";
import { Text0, Text1, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import {
  EmailInputBar,
  InputBar,
  PasswordInputBar,
  UsernameInputBar,
} from "../../../ui/InputBar";
import { ReqButton } from "../../../ui/ReqButton";
import { AuthState, RequestState, RequestTypes, SignupForm } from "../../../utils/types";
import { useDispatch, useSelector } from "react-redux";
import { AuthActions } from "../../../redux/slices/auth.slice";
import { useShake } from "../../../utils/hooks";
import { a } from "@react-spring/web";
import {
  validateEmail,
  validatePassword,
  validateUsername,
} from "../../../utils/regex";
import { RootState } from "../../../redux/store";
import { ReqActions } from "../../../redux/slices/req.slice";

export const Signup: React.FC = () => {
  const [signupForm, setSginupForm] = useState<SignupForm>({});
  const shakeAnimation = useShake(0, 2);
  const dispatch = useDispatch();
  const { isPending } = useSelector((state: RootState) => ({
    isPending: state.req.requestState === RequestState.Pending && state.req.reqType === RequestTypes.Signup,
  }));

  const handleSubmit = () => {
    if (
      !signupForm.email ||
      !signupForm.password ||
      !signupForm.username ||
      !validateEmail(signupForm.email) ||
      !validatePassword(signupForm.password) ||
      !validateUsername(signupForm.username)
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
    dispatch(ReqActions.setState({
      requestState: RequestState.Pending,
      reqType: RequestTypes.Signup
    }))
    dispatch(AuthActions.signup(signupForm));
  }
  return (
    <div className={classNames(styles.signupWrapper)}>
      <Text1
        text={strings.auth.signup}
        style={{ color: "var(--color-neutrals-n-500)" }}
      />
      <div className={classNames(styles.inputWrapper)}>
        <UsernameInputBar
          placeholder={strings.auth.enterUsername}
          value={signupForm.username}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSginupForm((prev) => ({ ...prev, username: e.target.value }))
          }
        />
        <EmailInputBar
          placeholder={strings.auth.enterEmail}
          value={signupForm.email}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSginupForm((prev) => ({ ...prev, email: e.target.value }))
          }
        />
        <PasswordInputBar
          placeholder={strings.auth.enterPass}
          value={signupForm.password}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSginupForm((prev) => ({ ...prev, password: e.target.value }))
          }
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
          text={strings.auth.signup}
          style={{ fontSize: "medium", width: "40%" }}
          onClick={handleSubmit}
        />
      </a.div>
    </div>
  );
};
