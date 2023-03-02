import error from '../modules/error';

function delElem(this: Elem & Containable, uid: Uid) {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (!this.hasElem(uid)) {
    throw error.failedToGet("delElem", uid);
  }

  if (this.container instanceof Map) {
    this.container.delete(uid);
    return;
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default delElem;