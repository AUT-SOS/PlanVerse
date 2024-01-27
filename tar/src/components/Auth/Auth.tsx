import React, { useState } from "react";
import styles from "./Auth.module.scss";
import { Background } from "../../ui/BackGround";
import strings from "../../utils/text";
import classNames from "classnames";
import { Text0, Text1, Text2, Text3 } from "../../ui/Text";
import { Login } from "./Login/Login";
import { DOC_ADDRESS } from "../../utils/consts";
import { Signup } from "./Signup/Signup";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { AuthState, Breakpoints } from "../../utils/types";
import { a, useTransition } from "@react-spring/web";
import { EmailValidation } from "./Signup/EmailValidation";
import { useBreakPoints } from "../../utils/hooks";
import { useNavigate } from "react-router-dom";
import { ReqButton } from "../../ui/ReqButton";
import { AuthActions } from "../../redux/slices/auth.slice";

export const Auth: React.FC = () => {
  const authState = useSelector((state: RootState) => state.auth.authState);
  const navigate = useNavigate();
  if (authState === AuthState.Authenticated){
    navigate("/home")
  }
  const [isLogin, setIsLogin] = useState<boolean>(
    window.location.pathname !== "/signup"
  );
  const transition = useTransition(authState, {
    from: {
      y: 200,
      opacity: "0",
    },
    enter: {
      y: 0,
      opacity: "1",
    },
    leave: {
      y: -200,
      opacity: "0",
    },
  });
  const transition2 = useTransition(isLogin, {
    from: {
      y: 100,
      opacity: "0",
    },
    enter: {
      y: 0,
      opacity: "1",
    },
    leave: {
      y: -100,
      opacity: "0",
    },
    config: {
      duration: 200,
    },
  });
  const breakPoint = useBreakPoints();
  return breakPoint > Breakpoints.Medium ? (
    <Background className={styles.authWrapper}>
      {transition((style, state) =>
        state === AuthState.Unauthenticated ? (
          <a.div style={style} className={styles.authContainer}>
            <AuthCover isLogin={isLogin} setIsLogin={setIsLogin} />
            <Login />
            <Signup />
          </a.div>
        ) : (
          state === AuthState.EmailValidate && (
            <a.div style={style} className={styles.otpContainer}>
              <EmailValidation />
            </a.div>
          )
        )
      )}
    </Background>
  ) : (
    <Background className={styles.authWrapper}>
      {transition((style, state) =>
        state === AuthState.Unauthenticated ? (
          <a.div style={style} className={styles.authContainer}>
            {transition2((style, state) =>
              state ? (
                <a.div style={style} className={styles.smallCont}>
                  <Login />
                </a.div>
              ) : (
                <a.div style={style} className={styles.smallCont}>
                  <Signup />
                </a.div>
              )
            )}
            
          </a.div>
        ) : (
          state === AuthState.EmailValidate && (
            <a.div style={style} className={styles.otpContainer}>
              <EmailValidation />
            </a.div>
          )
        )
      )}
      {authState === AuthState.Unauthenticated && <Text3
              text={!isLogin ? strings.auth.loginL : strings.auth.signupL}
              style={{
                cursor: "pointer",
                color: "var(--color-neutrals-n-400)",
                textDecoration: "underline",
              }}
              onClick={() => setIsLogin((prev) => !prev)}
            />}
    </Background>
  );
};

type Props = {
  isLogin: boolean;
  setIsLogin: (val: boolean) => void;
};

const AuthCover: React.FC<Props> = (props) => {
  return (
    <div
      className={classNames(styles.coverWrapper, {
        [styles.isRight]: props.isLogin,
      })}
    >
      {props.isLogin ? (
        <>
          <Text0
            style={{
              textAlign: "center",
              color: "var(--color-neutrals-on-primary)",
            }}
            text={strings.auth.loginWelcome}
          />
          <Text3
            text={strings.auth.signupL}
            style={{
              cursor: "pointer",
              color: "var(--color-link-light)",
              textDecoration: "underline",
            }}
            onClick={() => props.setIsLogin(!props.isLogin)}
          />
        </>
      ) : (
        <>
          <Text1
            style={{
              textAlign: "center",
              color: "var(--color-neutrals-on-primary)",
            }}
            text={strings.auth.signupWelcome}
          />
          <Text2
            text={strings.info}
            style={{
              textAlign: "center",
              color: "var(--color-neutrals-on-primary)",
            }}
          />
          <a href={DOC_ADDRESS} target="_blank">
            {strings.more}
          </a>
          <Text3
            text={strings.auth.loginL}
            style={{
              cursor: "pointer",
              color: "var(--color-link-light)",
              textDecoration: "underline",
            }}
            onClick={() => props.setIsLogin(!props.isLogin)}
          />
        </>
      )}
    </div>
  );
};
