import error from '../modules/error';

/* public */
function set(this: Elem, domElem: HTMLElement): HTMLElement {
  this.root = domElem;

  return this.root;
}

export default set;
