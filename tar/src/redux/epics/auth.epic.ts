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
import { toast } from "react-toastify";

export const loginEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.login.type),
    mergeMap((action) => {
      const loginInfo = action.payload as LoginForm;
      return API.login(loginInfo.email!, loginInfo.password!).pipe(
        mergeMap(() => {
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
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
      return API.signup(
        signupInfo.email!,
        signupInfo.password!,
        signupInfo.username!,
      ).pipe(
        mergeMap((res) => {
          document.cookie = `access_token=${res.responseHeaders.authorization}`
          
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.EmailValidate,
                exInfo: {
                  email: signupInfo.email
                }
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
        }),
        catchError(() => {
          showFailToastMessage("Email is already used");
          return of(ReqActions.setState({ requestState: RequestState.Error }));
        })
      );
    })
  );

  export const verificationEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.otpVerify.type),
    mergeMap((action) => {
      const otpPayload = action.payload as string;
      return API.otpVerify(
        otpPayload,
      ).pipe(
        mergeMap(() => {
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
        }),
        catchError(() => {
          return of(ReqActions.setState({ requestState: RequestState.Error }));
        })
      );
    })
  );