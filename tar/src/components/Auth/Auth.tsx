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

export const Auth: React.FC = () => {
  return (
    <div className={styles.authWrapper}>
      <div className={styles.authContainer}>
        <AuthCover />
        <Login/>
        <Signup/>
      </div>
    </div>
  );
};

const AuthCover: React.FC = (props) => {
  const [isLogin, setIsLogin] = useState<boolean>(
    window.location.pathname === "/login"
  );
  return (
    <div
      className={classNames(styles.coverWrapper, { [styles.isRight]: isLogin })}
    >
      {isLogin ? (
        <>
          <Text0 text={strings.auth.loginWelcome} />
          <Text3
            text={strings.auth.signupL}
            style={{ cursor: "pointer", color: "var(--color-link-light)", textDecoration: "underline"}}
            onClick={() => setIsLogin((prev) => !prev)}
          />
        </>
      ) : (
        <>
          <Text1 text={strings.auth.signupWelcome} />
          <Text2 text={strings.info} style={{textAlign: "center"}} />
          <a href={DOC_ADDRESS} target="_blank">{strings.more}</a>
          <Text3
            text={strings.auth.loginL}
            style={{ cursor: "pointer", color: "var(--color-link-light)", textDecoration: "underline"}}
            onClick={() => setIsLogin((prev) => !prev)}
          />
        </>
      )}
    </div>
  );
};
