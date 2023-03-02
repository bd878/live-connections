import error from '../modules/error';

function clearContainer(this: Elem & Containable): void {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (this.container instanceof Map) {
    this.container.clear();
  }

  throw error.wrongDataType(this.getName(), typeof this.container);
}

export default clearContainer;