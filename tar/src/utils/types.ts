export enum AuthState {
  Authenticated,
  Unauthenticated,
  EmailValidate,
}

export enum RequestState {
  None,
  Pending,
  Error,
}

export enum RequestTypes {
  Login,
  Signup,
  EmailValidate,
  ResendEmail,
  CreateProject,
  EditUser,
  JoinProject,
  ShareLink,
  EditProject,
}

export type JoinProjectType = {
  project_id: string;
  title: string;
  picture: string;
  members_number: number;
  members: { Username: string; ProfilePic: string }[];
};

export type LoginForm = {
  email?: string;
  password?: string;
};

export type SignupForm = {
  email?: string;
  password?: string;
  username?: string;
};

export type User = {
  id: string;
  email: string;
  username: string;
  profile_pic: string;
};

export type Project = {
  project_id: string;
  title: string;
  picture: string;
  description: string;
  owner_id: string;
  members_number: number;
};

export type SmallProject = {
  project_id: string;
  title: string;
  picture: string;
  members_number: string;
  is_admin: boolean;
};

export type Member = {
  id: string;
  username: string;
  email: string;
  profile_pic: string;
  is_admin: boolean;
};

export type CreateProject = {
  title: string;
  picture: string;
  description: string;
};

export type ShareLink = {
  id: string;
  emails: string[]
}

export type Task = {
  id: number;
  desc: number;
  memebers: User[];
  comments: TaskComment[];
};

export type TaskColumn = {
  id: number;
  name: string;
  tasks: Task[];
};

export type TaskComment = {
  id: number;
  sender: User;
  text: string;
};

export type UserEditType = {
  username: string;
  password: string;
  email: string;
  profile_pic: string;
};

export type IconProps = {
  color?: string;
  size: number;
  className?: string;
  bgColor?: string;
  secondColor?: string;
} & React.SVGProps<SVGSVGElement>;
export type IconComponent = React.FC<IconProps>;

export enum Breakpoints {
  XSmall = 300,
  Small = 500,
  Medium = 1000,
  Large = 1500,
}
