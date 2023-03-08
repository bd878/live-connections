import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import setId from '../../methods/setId';
import getName from '../../methods/getName';
import create from '../../methods/create';
import redraw from '../../methods/redrawBg';
import getUid from '../../methods/getUid';

class Cursor implements
  Elem,
  Creatable,
  Accessible,
  Moveable,
  Styleable,
  Redrawable,
  Identifable,
  Colored
{
  static cname: string = "Cursor";

  root: HTMLElement | null = null;
  id: string = '';
  CSSClass: string = "cursor";

  getName = getName;

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;
  move = move;

  redraw = redraw;

  setId = setId;
  getUid = getUid;
}

export default Cursor;
