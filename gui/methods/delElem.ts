import { isDeletable } from '../rtti';
import error from '../modules/error';

function delElem(this: Elem & Containable, uid: Uid): Elem & Containable {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (!this.hasElem(uid)) {
    throw error.failedToGet("delElem", uid);
  }

  if (this.container instanceof Map) {
    const elem = this.getElem(uid);
    ;(isDeletable(elem) && elem.free());
    this.container.delete(uid);
    return this;
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default delElem;