import React, { useState } from "react";
import styles from "./Signup.module.scss";
import classNames from "classnames";
import { Text0, Text1, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import { EmailInputBar, InputBar, PasswordInputBar } from "../../../ui/InputBar";
import { ReqButton } from "../../../ui/ReqButton";
import { SignupForm } from "../../../utils/types";



export const Signup: React.FC = () => {
  const [signupForm, setSginupForm] = useState<SignupForm>({});
  return (
    <div className={classNames(styles.signupWrapper)}>
      <Text1 text={strings.auth.signup} style={{color: "var(--color-neutrals-n-500)"}}/>
      <div className={classNames(styles.inputWrapper)}>
      <InputBar placeholder={strings.auth.enterUsername} value={signupForm.username} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSginupForm((prev) => ({...prev, username: e.target.value}))}/>
        <EmailInputBar placeholder={strings.auth.enterEmail} value={signupForm.email} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSginupForm((prev) => ({...prev, email: e.target.value}))}/>
        <PasswordInputBar title="Password should be 8 characters" placeholder={strings.auth.enterPass} value={signupForm.password} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSginupForm((prev) => ({...prev, password: e.target.value}))}/>
      </div>
      <ReqButton text={strings.auth.signup} style={{fontSize: "medium", width: "40%"}}/>
    </div>
  );
};
