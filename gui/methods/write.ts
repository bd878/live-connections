import error from '../modules/error';

/* public */
function write(this: Elem, content: string): void {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  this.root.textContent = content;
}

export default write;
