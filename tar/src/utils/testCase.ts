import { RANDOM_IMAGE } from "./consts";
import { Project, SmallProject, User } from "./types";

export const user1 : User = {
      id: "1",
      email: "",
      username: "Shayan",
      profilePic: RANDOM_IMAGE,
}

export const user2 : User = {
      id: "2",
      email: "",
      username: "Arshia",
      profilePic: RANDOM_IMAGE,
}

export const user3 : User = {
      id: "3",
      email: "",
      username: "Sara",
      profilePic: RANDOM_IMAGE,
}

export const user4 : User = {
      id: "4",
      email: "",
      username: "Sepehr",
      profilePic: RANDOM_IMAGE,
}

export const project1 : Project = {
      id: 1,
      name: "Narm 1",
      owner: user1,
      admins: [user1, user2],
      members: [user1, user2, user3],
      background: RANDOM_IMAGE,
      link: ""
}

export const sProj1 : SmallProject = {
      project_id: "1",
      title: "Planverse",
      picture: RANDOM_IMAGE,
      members_number: "5",
      is_admin: true,
}

export const sProj2 : SmallProject = {
      project_id: "2",
      title: "Test2",
      picture: "",
      members_number: "5",
      is_admin: true,
}