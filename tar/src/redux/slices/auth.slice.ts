import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { AuthState, LoginForm, SignupForm } from "../../utils/types";

type AuthSliceType = {
  authState: AuthState;
  exInfo?: any;
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
    },
    login(state, action: PayloadAction<LoginForm>) {
    },
    signup(state, action: PayloadAction<SignupForm>) {
    },
    
  },
});

export const AuthActions = AuthSlice.actions;
export const AuthReducers = AuthSlice.reducer;

export type LoginAction = ReturnType<typeof AuthActions.login>;
export type SignupAction = ReturnType<typeof AuthActions.signup>;
export type ChangeAuthStateAction = ReturnType<typeof AuthActions.changeAuthState>;

