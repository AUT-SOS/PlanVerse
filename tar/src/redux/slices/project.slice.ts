import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import {
  CreateProject,
  CreateStateType,
  CreateTaskType,
  JoinProjectType,
  Member,
  Project,
  ShareLink,
  SmallProject,
  State,
  Task,
} from "../../utils/types";

type ProjectSliceType = {
  myProjects: SmallProject[];
  fullProject?: Project;
  members?: Member[];
  joinProject?: JoinProjectType;
  states?: State[];
  task?: Task;
};

const initialState: ProjectSliceType = {
  myProjects: [],
};

const ProjectsSlice = createSlice({
  name: "Config",
  initialState,
  reducers: {
    setMyProjects(state, action: PayloadAction<SmallProject[]>) {
      state.myProjects = action.payload;
    },
    setFullProject(state, action: PayloadAction<Project>) {
      state.fullProject = action.payload;
    },
    setMembers(state, action: PayloadAction<Member[]>) {
      state.members = action.payload;
    },
    editMember(state, action: PayloadAction<Member>) {
      const index = state.members?.findIndex(
        (item) => item.id === action.payload.id
      );
      if (state.members && index) {
        state.members[index] = action.payload;
      }
    },
    setJoinProject(state, action: PayloadAction<JoinProjectType>) {
      state.joinProject = action.payload;
    },
    setStates(state, action: PayloadAction<State[]>) {
      state.states = action.payload;
    },
    setState(state, action: PayloadAction<State>) {
      const index = state.states?.findIndex(
        (item) => item.state_id === action.payload.state_id
      );
      if (state.states && index != undefined) {
        state.states[index] = action.payload;
      }
    },
    createTask(state, action: PayloadAction<CreateTaskType>) {},
    editTask(
      state,
      action: PayloadAction<CreateTaskType & { task_id: string }>
    ) {},
    createState(state, action: PayloadAction<CreateStateType>) {},
    editState(state, action: PayloadAction<State & { project_id: string }>) {
      const index = state.states?.findIndex(
        (item) => item.state_id === action.payload.state_id
      );
      if (state.states && index != undefined) {
        state.states[index] = action.payload;
      }
    },
    setTask(state, action: PayloadAction<Task | undefined>){
      state.task = action.payload
    }
  },
});

export const ProjectActions = {
  ...ProjectsSlice.actions,
  createProject: createAction<CreateProject>("Proj/CreateProject"),
  getMyProjects: createAction("Proj/GetMyProjects"),
  getFullProject: createAction<string>("Proj/GetFullProject"),
  changeMemberRole: createAction<{
    projectId: string;
    userId: string;
    isPromote: boolean;
  }>("Proj/ChangeMemberRole"),
  showProject: createAction<string>("Proj/ShowProject"),
  joinProject: createAction<string>("Proj/JoinProject"),
  shareLink: createAction<ShareLink>("Proj/ShareLink"),
  editProject: createAction<CreateProject & { id: string }>("Proj/EditProject"),
  deleteProject: createAction<{project_id: string, isDelete: boolean}>("Proj/DeleteProject"),
  getStates: createAction<string>("Proj/GetStates"),
  getState: createAction<{ projId: string; stateId: string }>("Proj/GetState"),
  deleteState: createAction<{ projId: string; stateId: string }>(
    "Proj/DeleteState"
  ),
  deleteTask: createAction<{ project_id: string; task_id: string }>(
    "Proj/DeleteTask"
  ),
  changeState: createAction<{
    project_id: string;
    task_id: string;
    state_id: string;
  }>("Proj/ChangeState"),
  getTask: createAction<string>("Proj/GetTask"),
  assign: createAction<{project_id: string, task_id: string, performer_id: string, isAdd: boolean}>("Proj/Assign")
};

export const ProjectReducer = ProjectsSlice.reducer;
