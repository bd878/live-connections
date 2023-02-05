import addCSSClass from './addCSSClass';
import { isStyleable } from '../rtti';

/* public */
function create(this: Elem): HTMLElement {
  this.root = document.createElement("div");
  if (isStyleable(this)) {
    addCSSClass.call(this, this.CSSClass);
  }
  return this.root;
}

export default create;
