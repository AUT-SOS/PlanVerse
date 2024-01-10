import { useSpring } from "@react-spring/web";

export const useShake = (
  shakeXstart: number,
  shakeXend: number,
) => {
  const [{ x }, api] = useSpring(() => ({
    from: { x: 0, y: 0 },
  }));
  const xInterpolate = x.to(
    [0, 0.25 , 0.55, 1],
    [shakeXstart, shakeXend, shakeXstart, shakeXend]
  );


  return {
    x: xInterpolate,
    api
  };
};
