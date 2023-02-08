import error from '../modules/error';

/* public */
function hasElem(this: Elem & Containable, key: string): boolean {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (this.container instanceof Map) {
    return this.container.has(key);
  }

  throw error.wrongDataType(this.name, typeof this.container);
}

export default hasElem;