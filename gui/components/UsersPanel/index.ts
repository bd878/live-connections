import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import write from '../../methods/write';

class UserPanel implements
  Elem,
  Creatable,
  Accessible,
  Writable
{
  root: HTMLElement | null = null;
  name: string = "user-panel";

  constructor() {}

  get = get;
  has = has;
  create = create;
  write = write;
}

const panel = new UserPanel();

export default panel;
