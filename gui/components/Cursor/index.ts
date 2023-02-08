import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import create from '../../methods/create';

class Cursor implements
  Elem,
  Creatable,
  Accessible,
  Moveable,
  Styleable
{
  root: HTMLElement | null = null;
  name: string = "cursor";
  CSSClass: string = "cursor";

  constructor() {}

  get = get;
  has = has;
  create = create;
  move = move;
}

export default Cursor;
