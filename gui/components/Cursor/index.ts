import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import create from '../../methods/create';
import redraw from './redraw';

class Cursor implements
  Elem,
  Creatable,
  Accessible,
  Moveable,
  Styleable,
  Redrawable
{
  root: HTMLElement | null = null;
  name: string = "cursor";
  CSSClass: string = "cursor";

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;
  move = move;

  redraw = redraw;
}

export default Cursor;
