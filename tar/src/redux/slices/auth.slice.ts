import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { AuthState, LoginForm, SignupForm } from "../../utils/types";

type authSliceType = {
  authState: AuthState;
};

const initialState: authSliceType = {
  authState: AuthState.Unauthenticated,
};

const AuthSlice = createSlice({
  name: "Auth",
  initialState,
  reducers: {
    changeAuthState(state, action: PayloadAction<AuthState>) {
      state.authState = action.payload;
    },
    login(state, action: PayloadAction<LoginForm>) {
      state.authState = AuthState.Pending;
    },
    signup(state, action: PayloadAction<SignupForm>) {
      state.authState = AuthState.Pending;
    },
    
  },
});

export const AuthActions = AuthSlice.actions;
export const AuthReducers = AuthSlice.reducer;

export type LoginAction = ReturnType<typeof AuthActions.login>;
export type SignupAction = ReturnType<typeof AuthActions.signup>;
export type ChangeAuthStateAction = ReturnType<typeof AuthActions.changeAuthState>;

