import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { AuthState, LoginForm, SignupForm, User } from "../../utils/types";

type UsersSliceType = {
  me?: User;
};

const initialState: UsersSliceType = {
  me: undefined,
};

const UsersSlice = createSlice({
  name: "Users",
  initialState,
  reducers: {
    setMe(state, action: PayloadAction<User>) {
      state.me = action.payload
    },
  },
});

export const UserActions = UsersSlice.actions;
export const UserReducers = UsersSlice.reducer;

