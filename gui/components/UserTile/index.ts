import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import redraw from './redraw';

class UserTile implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Redrawable
{
  root: HTMLElement | null = null;
  name: string = "user-tile";
  CSSClass: string = "user-tile";

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;

  redraw = redraw;
}

export default UserTile;
