import { ofType } from "redux-observable";
import { EMPTY, catchError, merge, mergeMap, of } from "rxjs";
import { API } from "../../api/API";
import { Epic, handleError } from "./epic";
import {
  Member,
  Project,
  RequestState,
  User,
  UserEditType,
} from "../../utils/types";
import { ProjectActions } from "../slices/project.slice";
import { UserActions } from "../slices/user.slice";
import { ReqActions } from "../slices/req.slice";

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
                of(
                  UserActions.setMe({
                    ...resObj,
                    profile_pic:
                      resObj.profile_pic.length > 0
                        ? resObj.profile_pic
                        : "/public//DefaultPFP.jpg",
                    id: myId,
                  } as User)
                ),
                of(ReqActions.setState({ requestState: RequestState.None }))
              );
            })
          );
        }),
        handleError("Password is incorrect")
      );
    })
  );
