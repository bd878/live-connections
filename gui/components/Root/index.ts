import error from '../../modules/error';
import append from '../../methods/append';
import get from '../../methods/get';
import has from '../../methods/has';
import set from '../../methods/set';

function setRoot(this: Root, domElem: HTMLElement): HTMLElement {
  const result = set.call(this, domElem);

  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  this.root.classList.add("root");

  return result;
}

class Root implements
  Elem,
  Appendable,
  Settable,
  Accessible
{
  root: HTMLElement | null = null;
  name: string = "root";

  constructor() {}

  append = append;
  set = setRoot;
  has = has;
  get = get;
}

const root = new Root();

export default root;
