import get from '../../methods/get';
import has from '../../methods/has';
import create from './create';

class Main implements
  Elem,
  Creatable,
  Accessible
{
  root: HTMLElement | null = null;
  name: string = "main";

  constructor() {}

  get = get;
  has = has;
  create = create;
}

const main = new Main();

export default main;
export { Main };
