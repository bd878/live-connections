
type AreaName = string;

type UserName = string;

type TextRecord = {
  value: string;
  updatedAt: number; // int64
  createdAt: number; // int64, id
};

type Uid = string;

type Id = UserName;

type Color = string;

type ABuffer = any;

type DView = any;

type Fn<A = any, R = any> = (args?: A) => R;

type CoordsEvent = {
  xPos: number;
  yPos: number;
  name: string;
};

type Coords = {
  xPos: number;
  yPos: number;
};

type BaseContainer = 
  | Record<string, any>
  | Map<string, any>;

interface Elem {
  root: HTMLElement | null;
  getName(): string;
}

interface Styleable {
  CSSClass: string;
}

type AuthOkEvent = {
  text: string;
};

type UsersOnlineEvent = {
  users: UserName[];
};

type TextInputEvent = {
  name: string;
  text: string;
};

type TitlesListEvent = {
  name: string;
  records: TextRecord[];
};

interface Appendable {
  append(I: Elem): void;
}

interface Settable {
  set(domElem: HTMLElement): HTMLElement;
}

interface Creatable {
  rootName?: string; // TODO: assign root div name to each component
  create(id: Id): HTMLElement;
}

interface Deletable {
  free(): boolean;
}

interface Writable {
  write(content: string): void;
}

interface Accessible {
  has(): boolean;
  get(): HTMLElement;
}

interface Moveable {
  move(x: number, y: number): void;
}

interface Colored {
  color: string;
}

interface Redrawable {
  redraw(...args: any[]): void;
}

interface Identifable {
  id: Id;
  setId(id: Id): void;
  getUid(): Uid;
}

interface Containable<C extends BaseContainer = BaseContainer> {
  container: C;

  hasElem(uid: Uid): boolean;
  getElem(uid: Uid): Elem;
  addElem(uid: Uid, elem: Elem): void;
  delElem(uid: Uid): void;
}

interface Clearable {
  clear(): void; /* virtual */
}
