import get from '../../methods/get';
import has from '../../methods/has';
import getName from '../../methods/getName';
import create from '../../methods/create';
import setId from '../../methods/setId';
import redraw from '../../methods/redrawBg';

class UserTile implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Redrawable,
  Identifable,
  Colored
{
  static cname: string = "UserTile";

  root: HTMLElement | null = null;
  id: string = '';
  CSSClass: string = "user-tile";

  getName = getName;

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;

  redraw = redraw;

  setId = setId;
}

export default UserTile;
