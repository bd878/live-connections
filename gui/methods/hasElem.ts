import error from '../modules/error';

/* public */
function hasElem(this: Elem & Containable, uid: Uid): boolean {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (this.container instanceof Map) {
    return this.container.has(uid);
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default hasElem;