import { ofType } from "redux-observable";
import { EMPTY, catchError, merge, mergeMap, of } from "rxjs";
import { API } from "../../api/API";
import { Epic, handleError } from "./epic";
import {
  AuthState,
  Member,
  Project,
  RequestState,
  User,
  UserEditType,
} from "../../utils/types";
import { ProjectActions } from "../slices/project.slice";
import { UserActions } from "../slices/user.slice";
import { ReqActions } from "../slices/req.slice";
import { AuthActions } from "../slices/auth.slice";
import { showFailToastMessage } from "../../main";

export const editUserEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(UserActions.editUserInfo.type),
    mergeMap((action) => {
      const myId = state$.value.auth.myId;
      const info = action.payload as UserEditType;
      return API.editUser(
        info.username,
        info.password,
        info.email,
        info.profile_pic
      ).pipe(
        mergeMap(() => {
          return API.getUser(myId!).pipe(
            mergeMap((res) => {
              const resObj = res.response as any;
              location.reload();
              return merge(
                of(ReqActions.setState({ requestState: RequestState.None }))
              );
            })
          );
        }),
        handleError("Password is incorrect")
      );
    })
  );

export const deleteUserEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(UserActions.deleteUser.type),
    mergeMap((action) => {
      const myId = state$.value.auth.myId!;
      return API.deleteUser(myId).pipe(
        mergeMap(() => {
          document.cookie = `access_token=logout`;
          location.reload();
          return EMPTY;
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );
