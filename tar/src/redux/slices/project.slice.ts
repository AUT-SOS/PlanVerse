import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import {
  CreateProject,
  JoinProjectType,
  Member,
  Project,
  ShareLink,
  SmallProject,
  State,
} from "../../utils/types";

type ProjectSliceType = {
  myProjects: SmallProject[];
  fullProject?: Project;
  members?: Member[];
  joinProject?: JoinProjectType;
  states?: State[];
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
    setStates(state, action: PayloadAction<State[]>){
      state.states = action.payload;
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
  deleteProject: createAction<string>("Proj/DeleteProject"),
  getStates: createAction<string>("Proj/GetStates"),
};

export const ProjectReducer = ProjectsSlice.reducer;
