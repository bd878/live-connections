import get from '../../methods/get';
import has from '../../methods/has';
import move from '../../methods/move';
import create from '../../methods/create';
import setId from '../../methods/setId';
import getName from '../../methods/getName';
import redraw from '../../methods/redrawBg';

class Square implements
  Elem,
  Creatable,
  Accessible,
  Moveable,
  Styleable,
  Redrawable,
  Identifable,
  Colored
{
  static cname: string = "Square";

  root: HTMLElement | null = null;
  id: string = '';
  CSSClass: string = "square";

  getName = getName;

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;
  move = move;

  redraw = redraw;

  setId = setId;
}

export default Square;
