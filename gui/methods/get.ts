import error from '../modules/error';

/* public */
function get(this: Elem): HTMLElement {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  return this.root;
}

export default get;
