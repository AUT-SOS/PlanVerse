import { StateObservable, combineEpics } from "redux-observable";
import { getMyIdEpic, loginEpic, resendEmailEpic, signupEpic, verificationEpic } from "./auth.epic";
import { Action, AnyAction } from "redux";
import { Observable, catchError, of } from "rxjs";
import { RootState } from "../store";
import { showFailToastMessage } from "../../main";
import { RequestState } from "../../utils/types";
import { ReqActions } from "../slices/req.slice";
import { changeMemberRoleEpic, createProjectEpic, deleteProject, editProjectEpic, getFullProject, getMyProjectsEpic, getStates, joinProjectEpic, shareLinkEpic, showProjectEpic } from "./project.epic";
import { deleteUserEpic, editUserEpic } from "./user.epic";

export interface Epic<
  Input extends Action = any,
  Output extends Action = any,
> {
  (
    action$: Observable<Input>,
    state$: StateObservable<RootState>,
  ): Observable<Output>;
}


export const handleError = <T>(message? : string) =>
  catchError<T, Observable<any>>((error) => {
    showFailToastMessage(message ?? error.message);
    return of(ReqActions.setState({ requestState: RequestState.Error }));
  });

export const rootEpic = combineEpics(
      signupEpic,
      loginEpic,
      verificationEpic,
      getMyIdEpic,
      resendEmailEpic,
      createProjectEpic,
      getMyProjectsEpic,
      getFullProject,
      changeMemberRoleEpic,
      editUserEpic,
      showProjectEpic,
      joinProjectEpic,
      shareLinkEpic,
      editProjectEpic,
      deleteUserEpic,
      deleteProject,
      getStates,
)