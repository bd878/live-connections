import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from '../../methods/create';

class Area implements
  Elem,
  Creatable,
  Accessible,
  Writable,
  Styleable
{
  root: HTMLElement | null = null;
  name: string = "area";
  CSSClass: string = "area";

  constructor() {}

  get = get;
  has = has;
  create = create;
  write = write;
}

const area = new Area();

export default area;
export { Area };
