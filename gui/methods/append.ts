import error from '../modules/error';

/* public */
function append(this: Elem, other: Elem) {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  if (!other.root) {
    throw error.noElementCreated(other.getName());
  }

  this.root.append(other.root);
}

export default append;
