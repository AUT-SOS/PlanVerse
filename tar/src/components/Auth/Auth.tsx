import React, { useState } from "react";
import styles from "./Auth.module.scss";
import { Background } from "../../ui/BackGround";
import { Title } from "../../ui/Title";
import strings from "../../utils/text";
import classNames from "classnames";
import { Text0, Text1, Text2, Text3 } from "../../ui/Text";
import { Login } from "./Login/Login";
import { DOC_ADDRESS } from "../../utils/consts";
import { Signup } from "./Signup/Signup";
import { useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { AuthState } from "../../utils/types";
import { a, useTransition } from "@react-spring/web";
import { EmailValidation } from "./Signup/EmailValidation";
import { useBreakPoints } from "../../utils/hooks";

export const Auth: React.FC = () => {
  const authState = useSelector((state: RootState) => state.auth.authState);
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
  
  return (
    <Background className={styles.authWrapper}>
      {transition((style, state) =>
        state === AuthState.Unauthenticated ? (
          <a.div style={style} className={styles.authContainer}>
            <AuthCover />
            <Login />
            <Signup />
          </a.div>
        ) : (
          state === AuthState.EmailValidate && (
            <a.div style={style} className={styles.otpContainer}>
              <EmailValidation/>
            </a.div>
          )
        )
      )}
    </Background>
  );
};

const AuthCover: React.FC = (props) => {
  const [isLogin, setIsLogin] = useState<boolean>(
    window.location.pathname === "/login"
  );
  const breakpoint = useBreakPoints();
  console.log(">>", breakpoint);
  
  return (
    <div
      className={classNames(styles.coverWrapper, { [styles.isRight]: isLogin })}
    >
      {isLogin ? (
        <>
          <Text0 style={{ textAlign: "center", color: "var(--color-neutrals-on-primary)" }} text={strings.auth.loginWelcome} />
          <Text3
            text={strings.auth.signupL}
            style={{
              cursor: "pointer",
              color: "var(--color-link-light)",
              textDecoration: "underline",
            }}
            onClick={() => setIsLogin((prev) => !prev)}
          />
        </>
      ) : (
        <>
          <Text1 style={{ textAlign: "center", color: "var(--color-neutrals-on-primary)" }}  text={strings.auth.signupWelcome} />
          <Text2 text={strings.info} style={{ textAlign: "center", color: "var(--color-neutrals-on-primary)" }} />
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
            onClick={() => setIsLogin((prev) => !prev)}
          />
        </>
      )}
    </div>
  );
};
