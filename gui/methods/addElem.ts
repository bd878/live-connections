import error from '../modules/error';

/* public */
function addElem(this: Elem & Containable, key: string, elem: Elem): boolean {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (this.container instanceof Map) {
    this.container.set(key, elem);
    return true;
  }

  throw error.wrongDataType(this.name, typeof this.container);
}

export default addElem;