import get from '../../methods/get';
import has from '../../methods/has';
import getName from '../../methods/getName';
import write from '../../methods/write';
import create from './create';

class UsersPanel implements
  Elem,
  Creatable,
  Accessible,
  Writable
{
  static cname: string = "UserPanel";

  root: HTMLElement | null = null;
  CSSClass: string = "users-panel";

  getName = getName;

  constructor() {}

  get = get;
  has = has;
  create = create;
  write = write;
}

const panel = new UsersPanel();

export default panel;
export { UsersPanel };
