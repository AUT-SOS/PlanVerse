import { ofType } from "redux-observable";
import { EMPTY, catchError, mergeMap, of } from "rxjs";
import { API } from "../../api/API";
import { Epic } from "./epic";
import { CreateStateType, CreateTaskType, State } from "../../utils/types";
import { ProjectActions } from "../slices/project.slice";
import { showFailToastMessage } from "../../main";

export const getStatesEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getStates.type),
    mergeMap((action) => {
      return API.Board.getStates(action.payload).pipe(
        mergeMap((res) => {
          const results = res.response as State[];
          results.sort((a, b) => Number(a.state_id) - Number(b.state_id))
          return of(ProjectActions.setStates(results));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const createTaskEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.createTask.type),
    mergeMap((action) => {
      const taskInfo = action.payload as CreateTaskType;
      return API.Board.createTask(
        taskInfo.project_id,
        taskInfo.state_id,
        taskInfo.title,
        taskInfo.back_ground_color,
        taskInfo.description
      ).pipe(
        mergeMap((res) => {
          return of(
            ProjectActions.getState({
              stateId: taskInfo.state_id,
              projId: taskInfo.project_id,
            })
          );
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const getStateEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getState.type),
    mergeMap((action) => {
      return API.Board.getState(
        action.payload.projId,
        action.payload.stateId
      ).pipe(
        mergeMap((res) => {
          return of(ProjectActions.setState(res.response as State));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const editTaskEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.editTask.type),
    mergeMap((action) => {
      const taskInfo = action.payload as CreateTaskType & { task_id: string };
      return API.Board.editTask(
        taskInfo.project_id,
        taskInfo.task_id,
        taskInfo.title,
        taskInfo.back_ground_color,
        taskInfo.description
      ).pipe(
        mergeMap(() => {
          return of(
            ProjectActions.getState({
              stateId: taskInfo.state_id,
              projId: taskInfo.project_id,
            })
          );
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

  export const createStateEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.createState.type),
    mergeMap((action) => {
      const stateInfo = action.payload as CreateStateType;
      return API.Board.createState(
        stateInfo.project_id,
        stateInfo.title,
        stateInfo.back_ground_color,
        stateInfo.admin_access
      ).pipe(
        mergeMap((res) => {
          return of(
            ProjectActions.getStates(stateInfo.project_id)
          );
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

  export const editStateEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.editState.type),
    mergeMap((action) => {
      const taskInfo = action.payload as State & { project_id: string };
      return API.Board.editState(
        taskInfo.project_id,
        taskInfo.state_id,
        taskInfo.title,
        taskInfo.back_ground_color,
        taskInfo.admin_access
      ).pipe(
        mergeMap(() => {
          return of(
            ProjectActions.getState({
              stateId: taskInfo.state_id,
              projId: taskInfo.project_id,
            })
          );
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );
