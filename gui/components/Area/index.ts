import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from '../../methods/create';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import getName from '../../methods/getName';
import addElemChild from '../../methods/addElemChild';
import delElemChild from '../../methods/delElemChild';
import redraw from './redraw';

class Area implements
  Elem,
  Creatable,
  Accessible,
  Writable,
  Styleable,
  Containable,
  Redrawable
{
  static cname: string = "Area";

  root: HTMLElement | null = null;
  CSSClass: string = "area";

  container: BaseContainer;

  getName = getName;

  constructor() {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;
  write = write;

  redraw = redraw;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElemChild;
  delElem = delElemChild;
}

const area = new Area();

export default area;
export { Area };
