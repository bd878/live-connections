import append from '../../methods/append';
import get from '../../methods/get';
import has from '../../methods/has';

class Root implements
  Elem,
  Appendable,
  Settable,
  Accessible,
  Styleable
{
  root: HTMLElement | null = null;
  name: string = "root";
  CSSClass: string = "root";

  constructor() {}

  append = append;
  set = setRoot;
  has = has;
  get = get;
}

const root = new Root();

export default root;
