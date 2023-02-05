import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';

class UsersList implements
  Elem,
  Creatable,
  Accessible,
  Styleable
{
  root: HTMLElement | null = null;
  name: string = "users-list";
  CSSClass: string = "users-list";

  constructor() {}

  get = get;
  has = has;
  create = create;
}

const list = new UsersList();

export default list;
export { UsersList };
