import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
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
  root: HTMLElement | null = null;
  name: string = "users-list";
  CSSClass: string = "users-list";

  container: BaseContainer;

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
}

const list = new UsersList();

export default list;
export { UsersList };
