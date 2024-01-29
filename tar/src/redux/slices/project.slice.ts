import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import { CreateProject, Member, Project, SmallProject } from "../../utils/types";


type ProjectSliceType = {
  myProjects: SmallProject[],
  fullProject?: Project,
  members?: Member[]

}

const initialState: ProjectSliceType = {
  myProjects: [],
};

const ProjectsSlice = createSlice({
  name: "Config",
  initialState,
  reducers: {
    setMyProjects(state, action: PayloadAction<SmallProject[]>) {
      state.myProjects = action.payload
    },
    setFullProject(state, action: PayloadAction<Project>){
      state.fullProject = action.payload
    },
    setMembers(state, action: PayloadAction<Member[]>){
      state.members = action.payload
    },
    editMember(state, action: PayloadAction<Member>){
      const index = state.members?.findIndex((item) => item.id === action.payload.id) 
      if (state.members && index) {
        state.members[index] = action.payload;
      } 
    }
  },
});

export const ProjectActions = {
  ...ProjectsSlice.actions,
  createProject: createAction<CreateProject>("Proj/CreateProject"),
  getMyProjects: createAction("Proj/GetMyProjects"),
  getFullProject: createAction<string>("Proj/GetFullProject"),
  changeMemberRole: createAction<{projectId: string, userId: string, isPromote: boolean}>("Proj/ChangeMemberRole"),
};

export const ProjectReducer = ProjectsSlice.reducer;
