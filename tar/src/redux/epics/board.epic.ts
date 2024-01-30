import { ofType } from "redux-observable";
import { EMPTY, catchError, iif, merge, mergeMap, of } from "rxjs";
import { API } from "../../api/API";
import { Epic } from "./epic";
import {
  CreateStateType,
  CreateTaskType,
  RequestState,
  RequestTypes,
  State,
  Task,
} from "../../utils/types";
import { ProjectActions } from "../slices/project.slice";
import { showFailToastMessage, showSuccessToastMessage } from "../../main";
import { ReqActions } from "../slices/req.slice";

export const getStatesEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getStates.type),
    mergeMap((action) => {
      return API.Board.getStates(action.payload).pipe(
        mergeMap((res) => {
          const results = res.response as State[];
          results.sort((a, b) => Number(a.state_id) - Number(b.state_id));
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
        taskInfo.description,
        taskInfo.index
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
          showFailToastMessage("Only admins can create tasks");
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
        taskInfo.description,
        taskInfo.index
      ).pipe(
        mergeMap(() => {
          if (
            state$.value.req.reqType === RequestTypes.EditTask &&
            state$.value.req.requestState === RequestState.Pending
          )
            showSuccessToastMessage("Editted successfully");
          return merge(
            of(ProjectActions.getStates(taskInfo.project_id)),
            of(ReqActions.setState({ requestState: RequestState.None }))
          );
        }),
        catchError(() => {
          showFailToastMessage("Only admins can edit tasks");
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
          return of(ProjectActions.getStates(stateInfo.project_id));
        }),
        catchError(() => {
          showFailToastMessage("Only admins can create states");
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
          showFailToastMessage("Only admins can edit states");
          return EMPTY;
        })
      );
    })
  );

export const deleteStateEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.deleteState.type),
    mergeMap((action) => {
      return API.Board.deleteState(
        action.payload.projId,
        action.payload.stateId
      ).pipe(
        mergeMap((res) => {
          return of(ProjectActions.getStates(action.payload.projId));
        }),
        catchError(() => {
          showFailToastMessage("Only admins can delete states");
          return EMPTY;
        })
      );
    })
  );

export const changeStateEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.changeState.type),
    mergeMap((action) => {
      const taskInfo = action.payload as {
        project_id: string;
        task_id: string;
        state_id: string;
      };
      return API.Board.changeState(
        taskInfo.project_id,
        taskInfo.task_id,
        taskInfo.state_id
      ).pipe(
        mergeMap(() => {
          return of(ProjectActions.getStates(taskInfo.project_id));
        }),
        catchError(() => {
          return EMPTY;
        })
      );
    })
  );

export const getTaskEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.getTask.type),
    mergeMap((action) => {
      return API.Board.getTask(action.payload).pipe(
        mergeMap((res) => {
          return of(ProjectActions.setTask(res.response as Task));
        }),
        catchError(() => {
          showFailToastMessage("There was an error");
          return EMPTY;
        })
      );
    })
  );

export const deleteTaskEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.deleteTask.type),
    mergeMap((action) => {
      console.log(">>", action.payload);

      return API.Board.deleteTask(
        action.payload.project_id,
        action.payload.task_id
      ).pipe(
        mergeMap(() => {
          return of(ProjectActions.getStates(action.payload.project_id));
        }),
        catchError(() => {
          showFailToastMessage("Only admins can delete tasks");
          return EMPTY;
        })
      );
    })
  );

export const assignEpic: Epic = (action$, state$) =>
  action$.pipe(
    ofType(ProjectActions.assign.type),
    mergeMap((action) => {
      return iif(
        () => action.payload.isAdd,
        API.Board.addAssign(
          action.payload.project_id,
          action.payload.task_id,
          action.payload.performer_id
        ),
        API.Board.removeAssign(
          action.payload.project_id,
          action.payload.task_id,
          action.payload.performer_id
        )
      ).pipe(
        mergeMap(() => {
          return of(ProjectActions.getStates(action.payload.project_id));
        }),
        catchError(() => {
          showFailToastMessage("Only admins can edit tasks");
          return EMPTY;
        })
      );
    })
  );
