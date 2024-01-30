import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import { AuthState, LoginForm, SignupForm, User, UserEditType } from "../../utils/types";

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
      state.me = action.payload;
    },
  },
});

export const UserActions = {
  ...UsersSlice.actions,
  editUserInfo: createAction<UserEditType>("Users/EditUser"),
  deleteUser: createAction("Users/Delete")
};
export const UserReducers = UsersSlice.reducer;
