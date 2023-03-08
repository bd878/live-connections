import error from '../modules/error';

function redrawBg(this: Elem & Colored) {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  this.root.style.backgroundColor = this.color;
}

export default redrawBg;
