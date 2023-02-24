import get from '../../methods/get';
import has from '../../methods/has';
import write from '../../methods/write';
import create from '../../methods/create';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
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

  redraw = redraw;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElemChild;
  delElem = delElemChild;
}

const area = new Area();

export default area;
export { Area };
