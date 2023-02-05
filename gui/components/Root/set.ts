import type { Root } from './index';
import error from '../../modules/error';
import set from '../../methods/set';

function setRoot(this: Root, domElem: HTMLElement): HTMLElement {
  const result = set.call(this, domElem);

  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  this.root.classList.add("root");

  return result;
}

export default setRoot;
