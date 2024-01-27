import { createAction } from "@reduxjs/toolkit";
import { CreateProject, SmallProject } from "../../utils/types";




export const ProjectActions = {
  createProject: createAction<CreateProject>("Proj/CreateProject"),
};
