import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from '../../methods/create';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import addElem from '../../methods/addElem';

class Area implements
  Elem,
  Creatable,
  Accessible,
  Writable,
  Styleable,
  Containable
{
  root: HTMLElement | null = null;
  name: string = "area";
  CSSClass: string = "area";

  container: BaseContainer;

  constructor() {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;
  write = write;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElem;
}

const area = new Area();

export default area;
export { Area };
