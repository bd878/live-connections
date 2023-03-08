import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import free from '../../methods/free';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import getName from '../../methods/getName';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import delElemChild from '../../methods/delElemChild';
import addElemChild from '../../methods/addElemChild';
import redraw from '../../methods/redrawBg';
import create from './create';

class Square implements
  Elem,
  Creatable,
  Accessible,
  Moveable,
  Styleable,
  Redrawable,
  Identifable,
  Colored,
  Containable
{
  static cname: string = "Square";

  root: HTMLElement | null = null;
  id: string = '';
  CSSClass: string = "square";

  container: BaseContainer;

  getName = getName;

  constructor(public color: string = '') {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;
  free = free;
  move = move;

  redraw = redraw;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElemChild;
  delElem = delElemChild;

  setId = setId;
  getUid = getUid;
}

export default Square;
