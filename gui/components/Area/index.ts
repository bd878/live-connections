import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from './create';

class Area implements
  Elem,
  Creatable,
  Accessible,
  Writable
{
  root: HTMLElement | null = null;
  name: string = "area";

  constructor() {}

  get = get;
  has = has;
  create = create;
  write = write;
}

const area = new Area();

export default area;
export { Area };
