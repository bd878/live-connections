import error from '../modules/error';

/* public */
function getElem(this: Elem & Containable, key: string): Elem {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (this.container instanceof Map) {
    return this.container.get(key);
  }

  throw error.wrongDataType(this.name, typeof this.container);
}

export default getElem;