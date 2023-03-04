import addCSSClass from './addCSSClass';
import getUid from '../misc/getUid';
import { isStyleable, isIdentifable } from '../rtti';

/* public */
function create(this: Elem): HTMLElement {
  this.root = document.createElement("div");
  ;(isStyleable(this) && addCSSClass.call(this, this.CSSClass));
  ;(isIdentifable(this) && (this.root.id = getUid(this.getName(), this.id)));
  return this.root;
}

export default create;
