import { PayloadAction, createAction, createSlice } from "@reduxjs/toolkit";
import { CreateProject, SmallProject } from "../../utils/types";


type ProjectSliceType = {
  myProjects: SmallProject[],

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
  },
});

export const ProjectActions = {
  ...ProjectsSlice.actions,
  createProject: createAction<CreateProject>("Proj/CreateProject"),
  getMyProjects: createAction("Proj/GetMyProjects"),

};

export const ProjectReducer = ProjectsSlice.reducer;
