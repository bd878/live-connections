import { create } from './static';
import error from '../../modules/error';
import getColorFromUserName from '../../misc/getColorFromUserName';

let users: Map<string, User> = new Map;
let _myName: UserName | null = null;
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

  users.set(name, user);
  _list.push(user);

  return user;
}

function flush() {
  users = new Map();
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

function getByName(name: UserName): User {
  const user = users.get(name);
  if (!user) throw error.failedToGet("users getByName", name);
  return user;
}

function setMyName(name: UserName) {
  _myName = name;
}

function myName(): UserName {
  if (!_myName) throw error.failedToGet("users myName");
  return _myName;
}

function me(): User {
  return getByName(myName());
}

export default {
  make,
  create,
  me,
  myName,
  setMyName,
  getByName,
  listNames,
  set,
  flush,
};
