import addCSSClass from './addCSSClass';
import getUid from '../misc/getUid';
import { isStyleable, isIdentifable } from '../rtti';

function create(this: Elem, id: Id = ''): HTMLElement {
  this.root = document.createElement("div");

  if (isStyleable(this)) {
    addCSSClass.call(this, this.CSSClass);
  }

  if (isIdentifable(this)) {
    this.setId(id);
    this.root.id = this.getUid();
  }

  return this.root;
}

export default create;
