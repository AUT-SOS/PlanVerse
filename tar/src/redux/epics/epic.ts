import { StateObservable, combineEpics } from "redux-observable";
import { loginEpic, signupEpic, verificationEpic } from "./auth.epic";
import { Action, AnyAction } from "redux";
import { Observable } from "rxjs";
import { RootState } from "../store";

export interface Epic<
  Input extends Action = any,
  Output extends Action = any,
> {
  (
    action$: Observable<Input>,
    state$: StateObservable<RootState>,
  ): Observable<Output>;
}



export const rootEpic = combineEpics(
      signupEpic,
      loginEpic,
      verificationEpic
)