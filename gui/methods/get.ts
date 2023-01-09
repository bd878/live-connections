import error from '../modules/error';

function get(this: Elem): HTMLElement {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  return this.root;
}

export default get;
