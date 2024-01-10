import { ofType } from "redux-observable";
import { AuthActions } from "../slices/auth.slice";
import { EMPTY, catchError, merge, mergeMap, of, timeout } from "rxjs";
import { API } from "../../api/API";
import { Epic } from "./epic";
import {
  AuthState,
  LoginForm,
  RequestState,
  SignupForm,
} from "../../utils/types";
import { showFailToastMessage } from "../../main";
import { ReqActions } from "../slices/req.slice";

export const loginEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.login.type),
    mergeMap((action) => {
      const loginInfo = action.payload as LoginForm;
      return API.login(loginInfo.email!, loginInfo.password!).pipe(
        mergeMap((res) => {
          /* const endReq = of(ReqActions.setState({requestState: RequestState.None}))*/
          return EMPTY;
        }),
        catchError(() => {
          showFailToastMessage("Invalid username or password");
          const endReq = of(
            ReqActions.setState({ requestState: RequestState.None })
          );
          return merge(
            endReq,
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Unauthenticated,
              })
            )
          );
        })
      );
    })
  );

export const signupEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.signup.type),
    mergeMap((action) => {
      const signupInfo = action.payload as SignupForm;
      return of(
        EMPTY
      ); /*API.signup(action.payload.email!, action.payload.password!, action.payload.username).pipe(
        mergeMap((res) => EMPTY),
        catchError(() => {
          toast.info("HI");
          return EMPTY;
        })
      );*/
    })
  );
