import error from '../modules/error';

function delElem(this: Elem & Containable, name: string) {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (!this.hasElem(name)) {
    throw error.failedToGet("delElem", name);
  }

  if (this.container instanceof Map) {
    this.container.delete(name);
    return;
  }

  throw error.wrongDataType(this.name, typeof this.container);
}

export default delElem;