import error from '../modules/error';

/* public */
function addElem(this: Elem & Containable, uid: Uid, elem: Elem): boolean {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (this.container instanceof Map) {
    this.container.set(uid, elem);
    return true;
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default addElem;