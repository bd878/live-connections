const users: User[] = [];

class User {
  constructor(
    public area: AreaName = '',
    public user: UserName = '',
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

  setToken(token: string): void {
    this.token = token;
  }

  define(areaName: AreaName, userName: UserName): void {
    ;(!this.area && (this.area = areaName));
    ;(!this.user && (this.name = userName));
  }
}

function makeUser(area: AreaName, name: UserName, token: string | null = null): User {
  const user = new User(area, name, token);
  users.push(user)
  return user;
}

export default User;
export { makeUser };
export { create } from './static';
