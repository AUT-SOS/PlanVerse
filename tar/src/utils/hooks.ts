import { useSpring } from "@react-spring/web";
import { RequestState, RequestTypes } from "./types";
import { useSelector } from "react-redux";
import { RootState } from "../redux/store";
import { useLayoutEffect, useState } from "react";

export const useShake = (shakeXstart: number, shakeXend: number) => {
  const [{ x }, api] = useSpring(() => ({
    from: { x: 0, y: 0 },
  }));
  const xInterpolate = x.to(
    [0, 0.25, 0.55, 1],
    [shakeXstart, shakeXend, shakeXstart, shakeXend]
  );

  return {
    x: xInterpolate,
    api,
  };
};

export const useRequestStates = (componentType: RequestTypes) => {
  const state = useSelector((state: RootState) => {
    const isMine = state.req.reqType === componentType;
    return {
      isPending: state.req.requestState === RequestState.Pending && isMine,
      errorState: state.req.requestState === RequestState.Error && isMine,
    };
  });
  return state
};

export const useBreakPoints = () => {
  const breakpoint = useSelector((state : RootState) => state.config.breakpoint) ?? window.innerWidth;
  return breakpoint
}
