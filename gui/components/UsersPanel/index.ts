import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from './create';

class UsersPanel implements
  Elem,
  Creatable,
  Accessible,
  Writable
{
  root: HTMLElement | null = null;
  name: string = "user-panel";
  CSSClass: string = "users-panel";

  constructor() {}

  get = get;
  has = has;
  create = create;
  write = write;
}

const panel = new UsersPanel();

export default panel;
export { UsersPanel };
