
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

interface Elem {
  root: HTMLElement | null;
  name: string;
}

type AuthOkEvent = {
  text: string;
};

type UsersOnlineEvent = {
  users: UserName[];
};

interface Appendable {
  append(I: Elem): void;
}

interface Settable {
  set(domElem: HTMLElement): HTMLElement;
}

interface Creatable {
  create(): HTMLElement;
}

interface Writable {
  write(content: string): void;
}

interface Accessible {
  has(): boolean;
  get(): HTMLElement;
}