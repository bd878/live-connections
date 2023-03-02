import append from '../../methods/append';
import get from '../../methods/get';
import has from '../../methods/has';
import getName from '../../methods/getName';
import setRoot from './set';

class Root implements
  Elem,
  Appendable,
  Settable,
  Accessible,
  Styleable
{
  static cname: string = "Root";

  root: HTMLElement | null = null;
  CSSClass: string = "root";

  getName = getName;

  constructor() {}

  append = append;
  set = setRoot;
  has = has;
  get = get;
}

const root = new Root();

export default root;
export { Root };
