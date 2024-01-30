import { ofType } from "redux-observable";
import { AuthActions } from "../slices/auth.slice";
import {
  catchError,
  merge,
  mergeMap,
  of,
} from "rxjs";
import { API } from "../../api/API";
import { Epic, handleError } from "./epic";
import {
  AuthState,
  LoginForm,
  RequestState,
  SignupForm,
  User,
} from "../../utils/types";
import { showFailToastMessage, showSuccessToastMessage } from "../../main";
import { ReqActions } from "../slices/req.slice";
import { UserActions } from "../slices/user.slice";

export const loginEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.login.type),
    mergeMap((action) => {
      const loginInfo = action.payload as LoginForm;
      return API.login(loginInfo.email!, loginInfo.password!).pipe(
        mergeMap((res) => {
          document.cookie = `access_token=${res.responseHeaders.authorization}`;
          const uid = JSON.parse(JSON.stringify(res.response)).user_id;
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
                myId: uid,
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None })),
            API.getUser(uid).pipe(
              mergeMap((res) => {
                const resObj = res.response as any;
                return of(
                  UserActions.setMe({
                    ...resObj,
                    profile_pic:
                      resObj.profile_pic.length > 0
                        ? resObj.profile_pic
                        : "/public//DefaultPFP.jpg",
                    id: uid,
                  } as User)
                );
              })
            )
          );
        }),
        handleError("Invalid username or password")
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
          document.cookie = `access_token=${res.responseHeaders.authorization}`;

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
        handleError("Email already exists")
      );
    })
  );

export const getMyIdEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(AuthActions.getMyUserId.type),
    mergeMap(() => {
      return API.getMyId().pipe(
        mergeMap((res) => {
          const myId = JSON.stringify(res.response);
          return merge(
            of(AuthActions.setMyUserId(JSON.stringify(res.response))),
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
              })
            ),
            API.getUser(myId).pipe(
              mergeMap((res) => {
                const resObj = res.response as any;
                
                return of(
                  UserActions.setMe({
                    ...resObj,
                    profile_pic:
                      resObj.profile_pic.length > 0
                        ? resObj.profile_pic
                        : "/public//DefaultPFP.jpg",
                    id: myId,
                  } as User)
                );
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
        mergeMap((res) => {
          const myId = JSON.stringify((res.response as any).user_id);
          return merge(
            of(
              AuthActions.changeAuthState({
                authState: AuthState.Authenticated,
              })
            ),
            of(ReqActions.setState({ requestState: RequestState.None })),
            API.getUser(myId).pipe(
              mergeMap((res) => {
                const resObj = res.response as any;
                return of(
                  UserActions.setMe({
                    ...resObj,
                    profile_pic:
                      resObj.profile_pic.length > 0
                        ? resObj.profile_pic
                        : "/public//DefaultPFP.jpg",
                    id: myId,
                  } as User)
                );
              })
            )
          );
        }),
        handleError("Wrong Code")
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
        handleError("There was an error resendig")
      );
    })
  );
