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
} from "../../utils/types";
import { ReqActions } from "../slices/req.slice";
import { ProjectActions } from "../slices/project.slice";
import { navigate } from "../../utils/configs";

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
            )
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
      const newMember = state$.value.project.members?.find(
        (item) => item.id === action.payload.userId
      )!;
      newMember.is_admin = !newMember.is_admin;
      return merge(
        of(ProjectActions.editMember(newMember)),
        iif(
          () => action.payload.isPromote,
          API.promote(action.payload.projectId, action.payload.userId),
          API.demote(action.payload.projectId, action.payload.userId)
        ).pipe(
          mergeMap(() =>
            API.getProjectMembers(action.payload).pipe(
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
          return of(ProjectActions.setJoinProject(res.response as JoinProjectType))
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
          return of(ReqActions.setState({requestState: RequestState.None}))
        }),
        catchError(() => {
          navigate(`/projects/${action.payload}`);
          return of(ReqActions.setState({requestState: RequestState.Error}));
        })
      );
    })
  );