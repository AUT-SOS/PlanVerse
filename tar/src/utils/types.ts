export enum AuthState  {
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
  CreateProject
}

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
  profilePic: string;
};

export type Project = {
  id: number;
  name: string;
  owner: User;
  admins: User[];
  members: User[];
  background: string;
  link: string;
};

export type SmallProject = {
  project_id: string,
  title: string,
  picture: string,
  members_number: string,
  is_admin: boolean
}

export type CreateProject = {
  title: string,
  picture: string,
  description: string,
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