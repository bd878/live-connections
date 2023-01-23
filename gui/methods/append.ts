import error from '../modules/error';

/* public */
function append(this: Elem, other: Elem) {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  if (!other.root) {
    throw error.noElementCreated(other.name);
  }

  this.root.append(other.root);
}

export default append;
