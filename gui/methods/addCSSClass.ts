import error from '../modules/error';

/* protected */
function addCSSClass(this: Elem & Styleable, className: string) {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  this.root.classList.add(className);
}

export default addCSSClass;
