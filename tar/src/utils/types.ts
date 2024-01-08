
export type User = {
  id: number;
  email: string;
  username: string;
  profile_pic: string;
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
