import error from '../modules/error';

/* public */
function addElem(this: Elem & Containable, uid: Uid, elem: Elem): Elem & Containable {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (this.container instanceof Map) {
    this.container.set(uid, elem);
    return this;
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default addElem;