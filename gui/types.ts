
type AreaName = string;

type UserName = string;

type ABuffer = any;

type DView = any;

type Fn<A = any, R = any> = (args?: A) => R;

type MouseMoveEvent = {
  xPos: number;
  yPos: number;
  name: string;
};

type AuthOkEvent = {
  text: string;
};

type UsersOnlineEvent = {
  users: UserName[];
};

interface Appendable {
  append(I: Element): void;
}