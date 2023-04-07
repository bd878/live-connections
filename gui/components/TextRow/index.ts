import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import getName from '../../methods/getName';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import addElem from '../../methods/addElem';
import getRecord from './getRecord';
import redraw from './redraw';

class TextRow implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Redrawable
{
  static cname: string = "TextRow";

  root: HTMLElement | null = null;
  CSSClass: string = "text-row";

  getName = getName;

  constructor(protected record: TextRecord) {
  }

  get = get;
  has = has;
  create = create;
  redraw = redraw;
  getRecord = getRecord;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElem;
}

export default TextRow;
