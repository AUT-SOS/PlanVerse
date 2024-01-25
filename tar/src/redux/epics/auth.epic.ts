import { ofType } from "redux-observable";
import { AuthActions } from "../slices/auth.slice";
import { EMPTY, catchError, merge, mergeMap, mergeMapTo, of, timeout } from "rxjs";
import { API } from "../../api/API";
import { Epic } from "./epic";
import {
  AuthState,
  LoginForm,
  RequestState,
  SignupForm,
} from "../../utils/types";
import { showFailToastMessage, showSuccessToastMessage } from "../../main";
import { ReqActions } from "../slices/req.slice";
import { toast } from "react-toastify";
import strings from "../../utils/text";

export const loginEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.login.type),
    mergeMap((action) => {
      const loginInfo = action.payload as LoginForm;
      return API.login(loginInfo.email!, loginInfo.password!).pipe(
        mergeMap((res) => {
          document.cookie = `access_token=${res.responseHeaders.authorization}; expires=Tue, 19 Jan 2038 04:14:07 GMT"`;
          const uid = JSON.parse(JSON.stringify(res.response)).user_id;
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
                myId: uid

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
        signupInfo.username!
      ).pipe(
        mergeMap((res) => {
          document.cookie = `access_token=${res.responseHeaders.authorization}; expires=Tue, 19 Jan 2038 04:14:07 GMT`;

          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.EmailValidate,
                exInfo: {
                  email: signupInfo.email,
                },
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
        }),
        catchError((error) => {
          showFailToastMessage(error.message);
          return of(ReqActions.setState({ requestState: RequestState.Error }));
        })
      );
    })
  );

export const getMyIdEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.getMyUserId.type),
    mergeMap(() => {
      return API.getMyId().pipe(
        mergeMap((res) => {
          return merge(
            of(AuthActions.setMyUserId(JSON.stringify(res.response))),
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
              })
            )
          );
        }),
        catchError(() => {
          return of(
            AuthActions.changeAuthState({
              authState: AuthState.Unauthenticated,
            })
          );
        })
      );
    })
  );

export const verificationEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.otpVerify.type),
    mergeMap((action) => {
      const otpPayload = action.payload as string;
      return API.otpVerify(otpPayload).pipe(
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

  export const resendEmailEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.resendEmail.type),
    mergeMap(() => {
      return API.resendEmail().pipe(
        mergeMap(() => {
          showSuccessToastMessage("Email was resent");
          return merge(
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
        }),
        catchError((error) => {
          showFailToastMessage(error.message);
          return of(ReqActions.setState({ requestState: RequestState.Error }));
        })
      );
    })
  );
