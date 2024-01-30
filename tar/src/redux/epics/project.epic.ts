import { ofType } from "redux-observable";
import { AuthActions } from "../slices/auth.slice";
import { EMPTY, catchError, iif, merge, mergeMap, mergeMapTo, of } from "rxjs";
import { API } from "../../api/API";
import { Epic, handleError } from "./epic";
import {
  AuthState,
  CreateProject,
  JoinProjectType,
  Member,
  Project,
  RequestState,
  SignupForm,
  SmallProject,
  State,
} from "../../utils/types";
import { ReqActions } from "../slices/req.slice";
import { ProjectActions } from "../slices/project.slice";
import { navigate } from "../../utils/configs";
import { showFailToastMessage, showSuccessToastMessage } from "../../main";

export const createProjectEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.createProject.type),
    mergeMap((action) => {
      const createProjectInfo = action.payload as CreateProject;
      return API.createProject(
        createProjectInfo.title,
        createProjectInfo.description,
        createProjectInfo.picture
      ).pipe(
        mergeMap((res) => {
          location.reload();
          return of(ReqActions.setState({ requestState: RequestState.None }));
        }),
        handleError()
      );
    })
  );

export const getMyProjectsEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getMyProjects.type),
    mergeMap(() => {
      return API.getMyProjects().pipe(
        mergeMap((res) => {
          return of(
            ProjectActions.setMyProjects(res.response as SmallProject[])
          );
        }),
        handleError()
      );
    })
  );

export const getFullProject: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getFullProject.type),
    mergeMap((action) => {
      return API.getFullProject(action.payload).pipe(
        mergeMap((res) => {
          return merge(
            of(ProjectActions.setFullProject(res.response as Project)),
            API.getProjectMembers(action.payload).pipe(
              mergeMap((res) => {
                return of(ProjectActions.setMembers(res.response as Member[]));
              })
            ),
            of(ProjectActions.getStates(action.payload))
          );
        }),
        catchError(() => {
          navigate("/home");
          return EMPTY;
        })
      );
    })
  );

export const changeMemberRoleEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.changeMemberRole.type),
    mergeMap((action) => {
      const newMember = JSON.parse(
        JSON.stringify(
          state$.value.project.members?.find(
            (item) => item.id == action.payload.userId
          )!
        )
      );
      newMember.is_admin = !newMember.is_admin;
      return merge(
        of(ProjectActions.editMember(newMember)),
        iif(
          () => action.payload.isPromote,
          API.promote(action.payload.projectId, action.payload.userId),
          API.demote(action.payload.projectId, action.payload.userId)
        ).pipe(
          mergeMap(() =>
            API.getProjectMembers(action.payload.projectId).pipe(
              mergeMap((res) => {
                return of(ProjectActions.setMembers(res.response as Member[]));
              })
            )
          ),
          handleError("Permission denied")
        )
      );
    })
  );

export const showProjectEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.showProject.type),
    mergeMap((action) => {
      return API.showProject(action.payload).pipe(
        mergeMap((res) => {
          return of(
            ProjectActions.setJoinProject(res.response as JoinProjectType)
          );
        }),
        catchError(() => {
          return EMPTY;
        })
      );
    })
  );

export const joinProjectEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.joinProject.type),
    mergeMap((action) => {
      return API.joinProject(action.payload).pipe(
        mergeMap(() => {
          navigate(`/projects/${action.payload}`);
          return of(ReqActions.setState({ requestState: RequestState.None }));
        }),
        catchError(() => {
          navigate(`/projects/${action.payload}`);
          return of(ReqActions.setState({ requestState: RequestState.Error }));
        })
      );
    })
  );

export const shareLinkEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.shareLink.type),
    mergeMap((action) => {
      return API.shareLink(action.payload.id, action.payload.emails).pipe(
        mergeMap((res) => {
          showSuccessToastMessage("Invite links were emailed successfully");
          return of(ReqActions.setState({ requestState: RequestState.None }));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const editProjectEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.editProject.type),
    mergeMap((action) => {
      const editProjectInfo = action.payload as CreateProject & { id: string };
      return API.editProject(
        editProjectInfo.id,
        editProjectInfo.title,
        editProjectInfo.description,
        editProjectInfo.picture
      ).pipe(
        mergeMap((rs) => {
          showSuccessToastMessage("Project was editted successfully");
          return merge(
            of(ReqActions.setState({ requestState: RequestState.None })),
            of(ProjectActions.getFullProject(editProjectInfo.id))
          );
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const deleteProject: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.deleteProject.type),
    mergeMap((action) => {
      return API.deleteProject(action.payload).pipe(
        mergeMap(() => {
          navigate("/home");
          return of(ReqActions.setState({ requestState: RequestState.None }));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const getStates: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getStates.type),
    mergeMap((action) => {
      return API.Board.getStates(action.payload).pipe(
        mergeMap((res) => {
          return of(ProjectActions.setStates(res.response as State[]));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );
