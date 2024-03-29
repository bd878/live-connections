import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from '../../methods/create';
import hasElem from '../../methods/hasElem';
import getName from '../../methods/getName';
import addChild from '../../methods/addChild';
import addElem from '../../methods/addElem';
import delElemChild from '../../methods/delElemChild';
import getElem from './getElem';
import free from './free';

class Area implements
  Elem,
  Creatable,
  Deletable,
  Accessible,
  Writable,
  Styleable,
  Containable
{
  static cname: string = "Area";

  root: HTMLElement | null = null;
  rootName: string | undefined;
  CSSClass: string = "area";

  container: BaseContainer;

  getName = getName;

  constructor() {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;
  free = free;
  write = write;

  addChild = addChild;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElem;
  delElem = delElemChild;
}

const area = new Area();

export default area;
export { Area };
