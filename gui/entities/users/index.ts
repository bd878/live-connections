import { create } from './static';
import getColorFromUserName from '../../misc/getColorFromUserName';

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

function listNames(): UserName[] {
  const names: UserName[] = [];

  for (let i = 0; i < _list.length; i++) {
    names.push(_list[i].name);
  }

  return names;
}

function set(area: AreaName, list: string[]) {
  flush();
  for (let i = 0; i < list.length; i++) {
    const name = list[i];
    make(area, name, getColorFromUserName(name));
  }
}

export default {
  make,
  create,
  users,
  listNames,
  set,
  flush,
};
