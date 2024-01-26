import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import { AuthState, LoginForm, SignupForm } from "../../utils/types";

type AuthSliceType = {
  authState: AuthState;
  exInfo?: any;
  myId?: string;
};

const initialState: AuthSliceType = {
  authState: AuthState.Unauthenticated,
};

const AuthSlice = createSlice({
  name: "Auth",
  initialState,
  reducers: {
    changeAuthState(state, action: PayloadAction<AuthSliceType>) {
      state.authState = action.payload.authState;
      state.exInfo = action.payload.exInfo;
      state.myId = action.payload.myId;
    },
    login(_state, _action: PayloadAction<LoginForm>) {},
    signup(state, _action: PayloadAction<SignupForm>) {
      state.exInfo = { email: _action.payload.email };
    },
    setMyUserId(state, action: PayloadAction<string>) {
      state.myId = action.payload;
    },
  },
});

export const AuthActions = {
  ...AuthSlice.actions,
  getMyUserId: createAction("Auth/GetMyId"),
  otpVerify: createAction<string>("Auth/OtpVerify"),
  resendEmail: createAction("Auth/ResendEmail")
};
export const AuthReducers = AuthSlice.reducer;

export type LoginAction = ReturnType<typeof AuthActions.login>;
export type SignupAction = ReturnType<typeof AuthActions.signup>;

