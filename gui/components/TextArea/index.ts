import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import getName from '../../methods/getName';
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
  textarea: HTMLElement | null = null;
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
}

export default TextArea;
