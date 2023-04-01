import get from '../../methods/get';
import has from '../../methods/has';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import getName from '../../methods/getName';
import create from '../../methods/create';
import turnReadonly from './turnReadonly';
import redraw from './redraw';

class TextArea implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Identifable,
  Redrawable
{
  static cname: string = "TextArea";

  root: HTMLElement | null = null;
  rootName: string | undefined = "textarea";
  id: string = '';
  CSSClass: string = "textarea";

  getName = getName;

  constructor() {}

  get = get;
  has = has;
  create = create;

  redraw = redraw;

  setId = setId;
  getUid = getUid;

  turnReadonly = turnReadonly;
}

export default TextArea;
