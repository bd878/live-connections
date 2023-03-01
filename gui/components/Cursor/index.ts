import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import create from '../../methods/create';
import redraw from '../../methods/redrawBg';

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
  root: HTMLElement | null = null;
  id: string = '';
  name: string = "cursor";
  CSSClass: string = "cursor";

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
