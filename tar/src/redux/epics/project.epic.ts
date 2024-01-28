import { ofType } from "redux-observable";
import { AuthActions } from "../slices/auth.slice";
import { merge, mergeMap, of } from "rxjs";
import { API } from "../../api/API";
import { Epic, handleError } from "./epic";
import {
  AuthState,
  CreateProject,
  RequestState,
  SignupForm,
  SmallProject,
} from "../../utils/types";
import { ReqActions } from "../slices/req.slice";
import { ProjectActions } from "../slices/project.slice";
import { showSuccessToastMessage } from "../../main";

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
    mergeMap((action) => {
      return API.getMyProjects(
      ).pipe(
        mergeMap((res) => {            
          return of(ProjectActions.setMyProjects(res.response as any));
        }),
        handleError()
      );
    })
  );
