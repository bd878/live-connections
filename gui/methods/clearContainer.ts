import error from '../modules/error';

function clearContainer(this: Elem & Containable): void {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (this.container instanceof Map) {
    this.container.clear();
  }

  throw error.wrongDataType(this.name, typeof this.container);
}

export default clearContainer;