import error from '../modules/error';

function set(this: Elem, domElem: HTMLElement): HTMLElement {
  this.root = domElem;

  return this.root;
}

export default set;
