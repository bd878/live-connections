
type AreaName = string;

type UserName = string;

type Fn<A = any, R = any> = (args: A) => R;

type MouseMoveEvent = {
  xPos: number;
  yPos: number;
  name: string;
};

type AuthOkEvent = {
  text: string;
};

type UsersOnlineEvent = {
  users: string[];
};
