import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import getName from '../../methods/getName';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import delElemChild from '../../methods/delElemChild';
import addElemChild from '../../methods/addElemChild';
import clearContainer from '../../methods/clearContainer';

class UsersList implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Containable,
  Clearable
{
  static cname: string = "UsersList";

  root: HTMLElement | null = null;
  CSSClass: string = "users-list";

  container: BaseContainer;

  getName = getName;

  constructor() {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;

  clear = clearContainer;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElemChild;
  delElem = delElemChild;
}

const list = new UsersList();

export default list;
export { UsersList };
