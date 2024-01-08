import { AnyAction, configureStore } from "@reduxjs/toolkit";
import { AuthReducers } from "./slices/auth.slice";
import { createEpicMiddleware } from "redux-observable";
import { rootEpic } from "./epics/epic";
import { combineReducers } from "redux";

const epicMiddleware = createEpicMiddleware<
  AnyAction,
  AnyAction,
  RootState,
  void
>();
const reducer = combineReducers({auth: AuthReducers});
export type RootState = ReturnType<typeof reducer>;
export const store = configureStore({
  reducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({ thunk: false }).concat(epicMiddleware as any),
});
epicMiddleware.run(rootEpic);
