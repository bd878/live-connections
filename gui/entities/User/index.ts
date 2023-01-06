import { create } from './static';

const users: User[] = [];

class User {
  constructor(
    public area: AreaName = '',
    public name: UserName = '',
    public token: string | null = null
  ) {
    this.isAuthed = this.isAuthed.bind(this);
    this.isNotAuthed = this.isNotAuthed.bind(this);
  }

  isAuthed(): boolean {
    return !!(this.area && this.name && this.token);
  }

  isNotAuthed(): boolean {
    return !this.isAuthed();
  }

  setToken(token: string) {
    this.token = token;
  }

  define(areaName: AreaName, userName: UserName) {
    ;(!this.area && (this.area = areaName));
    ;(!this.name && (this.name = userName));
  }
}

function make(area: AreaName, name: UserName, token: string | null = null): User {
  const user = new User(area, name, token);
  users.push(user)
  return user;
}

export default {
  make,
  create
};
