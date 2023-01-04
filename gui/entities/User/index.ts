class User {
  constructor(areaName = '', userName = '', token = '') {
    this.area = areaName;
    this.name = userName
    this.token = null;

    this.isAuthed = this.isAuthed.bind(this);
    this.isNotAuthed = this.isNotAuthed.bind(this);
  }

  isAuthed() {
    return !!(this.area && this.name && this.token);
  }

  isNotAuthed() {
    return !this.isAuthed();
  }

  setToken(token) {
    this.token = token;
  }

  define(areaName, userName) {
    ;(!this.area && (this.area = areaName));
    ;(!this.user && (this.name = userName));
  }
}

export default User;
export { create } from './static';
