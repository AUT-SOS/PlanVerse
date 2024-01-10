import React, { useCallback, useEffect, useRef, useState } from "react";
import styles from "./EmailValidation.module.scss";
import { Text2, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import { useSelector } from "react-redux";
import { RootState } from "../../../redux/store";
import { DigitInput } from "../../../ui/InputBar";
import { ReqButton } from "../../../ui/ReqButton";
import { AuthState } from "../../../utils/types";

export const EmailValidation: React.FC = (props) => {
  const [otp, setOtp] = useState<string[]>(new Array(5).fill(""));
  const { email, isPending } = useSelector((state: RootState) => ({
    email: state.auth.exInfo?.email,
    isPending: false,
  }));
  return (
    <div className={styles.otpContentWrapper}>
      <img
        className={styles.otpImg}
        src="../../../../public/Email.png"
        alt=""
      />
      <Text3
        text={strings.auth.otp_info}
        style={{ color: "var(--color-neutrals-n-500)", textAlign: "center" }}
      />
      <Text2
        text={email}
        style={{ color: "var(--color-neutrals-n-500)", textAlign: "center" }}
      />
      <OTPInput otp={otp} setOtp={setOtp} />
      <ReqButton
        isPending={isPending}
        text="Validate"
        style={{ width: "40%" }}
      />
    </div>
  );
};

type OTPInputProps = {
  otp: string[];
  setOtp: (otp: string[]) => void;
};

let currentOtpIndex = 0;
const OTPInput: React.FC<OTPInputProps> = (props) => {
  const [activeIndex, setActiveIndex] = useState<number>(0);
  const handleChange = useCallback(
    ({ target }: React.ChangeEvent<HTMLInputElement>) => {
      const { value } = target;
      const newOtp: string[] = [...props.otp];
      newOtp[currentOtpIndex] = value.substring(value.length - 1);
      if (!value) {
        setActiveIndex(currentOtpIndex - 1);
      } else {
        setActiveIndex(currentOtpIndex + 1);
      }
      props.setOtp(newOtp);
    },
    [props.otp, props.setOtp, currentOtpIndex]
  );

  const handleOnKeyDown = useCallback(
    ({ key }: React.KeyboardEvent<HTMLInputElement>, index: number) => {
      currentOtpIndex = index;
      if (key === "Backspace") {
        setActiveIndex(currentOtpIndex - 1);
      }
    },
    [props.otp, props.setOtp, currentOtpIndex]
  );

  const inputRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    inputRef.current?.focus();
  }, [activeIndex]);

  return (
    <div className={styles.otpInputWrapper}>
      {props.otp.map((_, index) => {
        return (
          <DigitInput
            ref={index === activeIndex ? inputRef : null}
            value={props.otp[index]}
            onChangeInput={handleChange}
            onKeyDown={(e) => handleOnKeyDown(e, index)}
          />
        );
      })}
    </div>
  );
};
