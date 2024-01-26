import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import styles from "./EmailValidation.module.scss";
import { Text2, Text3 } from "../../../ui/Text";
import strings from "../../../utils/text";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../../redux/store";
import { DigitInput } from "../../../ui/InputBar";
import { ReqButton, ResendEmailButton } from "../../../ui/ReqButton";
import { AuthState, RequestState, RequestTypes } from "../../../utils/types";
import { AuthActions } from "../../../redux/slices/auth.slice";
import { useRequestStates, useShake } from "../../../utils/hooks";
import { ReqActions } from "../../../redux/slices/req.slice";
import { OTPIcon } from "../../../ui/Icons/OTP";
import { secToString } from "../../../utils/configs";

export const EmailValidation: React.FC = (props) => {
  const [otp, setOtp] = useState<string[]>(new Array(5).fill(""));
  const { email } = useSelector((state: RootState) => ({
    email: state.auth.exInfo?.email,
  }));
  const shakeAnimation = useShake(0, 2);

  const { isPending, errorState } = useRequestStates(
    RequestTypes.EmailValidate
  );

  const dispatch = useDispatch();

  const handleValidation = useCallback(() => {
    const otpJoined = otp.join("");
    if (otpJoined.length < 5) {
      shakeAnimation.api.start({
        from: {
          x: 1.5,
        },
        to: {
          x: 0,
        },
      });
      dispatch(
        ReqActions.setState({
          requestState: RequestState.Error,
          reqType: RequestTypes.EmailValidate,
        })
      );
      return;
    }
    dispatch(AuthActions.otpVerify(otpJoined));
  }, [otp]);
  return (
    <div className={styles.otpContentWrapper}>
      <OTPIcon color="var(--color-neutrals-n-500)" size={60} />
      <Text3
        text={strings.auth.otp_info}
        style={{ color: "var(--color-neutrals-n-500)", textAlign: "center" }}
      />
      <Text2
        text={email}
        style={{ color: "var(--color-neutrals-n-500)", textAlign: "center" }}
      />
      <OTPInput otp={otp} setOtp={setOtp} error={errorState} />
      <ReqButton
        onClick={handleValidation}
        isPending={isPending}
        text="Validate"
        style={{ width: "40%", x: shakeAnimation.x }}
      />
      <div className={styles.ResendWrapper}>
        <Text3
          text="Did'nt receive email?"
          style={{ color: "var(--color-neutrals-n-500)", textAlign: "center" }}
        />
        <ResendEmail />
      </div>
    </div>
  );
};

type OTPInputProps = {
  otp: string[];
  setOtp: (otp: string[]) => void;
  error?: boolean;
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
            error={props.error}
          />
        );
      })}
    </div>
  );
};

const ResendEmail: React.FC = () => {
  const [time, setTime] = useState(10);
  const dispatch = useDispatch();
  useEffect(() => {
    if (time <= 0) {
      return
    }
    const timer = setTimeout(() => {
      setTime(time - 1);
    }, 1000);
    () => clearTimeout(timer);
  }, [time]);
  const text = useMemo(() => time <= 0 ? "Resend" : secToString(time), [time]);
  const handleClick = () => {
    dispatch(AuthActions.resendEmail());
    setTime(10);

  }
  return <ResendEmailButton onClick={handleClick} text={text} disable={time > 0} />;
};
