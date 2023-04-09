import get from '../../methods/get';
import has from '../../methods/has';
import getName from '../../methods/getName';
import create from '../../methods/create';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';

class Button implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Identifable
{
  static cname: string = "Button";

  root: HTMLElement | null = null;
  id: string = '';
  rootName: string | undefined = "button";
  CSSClass: string = "button";

  getName = getName;

  constructor() {
  }

  get = get;
  has = has;
  create = create;

  setId = setId;
  getUid = getUid;
}

export default Button;
