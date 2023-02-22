import { create } from './static';

let users: Record<string, User> = {};
let _list: User[] = []; // fast iterate

class User {
  constructor(
    public area: AreaName = '',
    public name: UserName = '',
    public color: Color = '',
    public token: string | null = null
  ) {}
}

function make(area: AreaName, name: UserName, color: Color, token: string | null = null): User {
  const user = new User(area, name, color, token);

  users[name] = user;
  _list.push(user);

  return user;
}

function flush() {
  users = {};
  _list = [];
}

function list(): { names: UserName[]; colors: Color[]; } {
  const names: UserName[] = [];
  const colors: Color[] = [];

  for (let i = 0; i < _list.length; i++) {
    const user = _list[i];

    names.push(user.name);
    colors.push(user.color);
  }

  return { names, colors };
}

export default {
  make,
  create,
  users,
  list,
  flush,
};
