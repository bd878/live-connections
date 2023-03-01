import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import redraw from '../../methods/redrawBg';

class UserTile implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Redrawable,
  Colored
{
  root: HTMLElement | null = null;
  id: string = '';
  name: string = "user-tile";
  CSSClass: string = "user-tile";

  constructor(public color: string = '') {}

  get = get;
  has = has;
  create = create;

  redraw = redraw;

  setId = setId;
  getUid = getUid;
}

export default UserTile;
