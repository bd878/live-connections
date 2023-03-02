import get from '../../methods/get';
import has from '../../methods/has';
import getName from '../../methods/getName';
import create from './create';

class Main implements
  Elem,
  Creatable,
  Accessible,
  Styleable
{
  static cname: string = "Main";

  root: HTMLElement | null = null;
  CSSClass: string = "main";

  getName = getName;

  constructor() {}

  get = get;
  has = has;
  create = create;
}

const main = new Main();

export default main;
export { Main };
