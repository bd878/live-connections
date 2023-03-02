import error from '../modules/error';

/* public */
function getElem(this: Elem & Containable, uid: Uid): Elem {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (this.container instanceof Map) {
    return this.container.get(uid);
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default getElem;