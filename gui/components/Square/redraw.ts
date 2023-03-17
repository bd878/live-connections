import type Square from './index';
import error from '../../modules/error';

function redraw(this: Square) {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  this.root.style.border = `2px solid ${this.color}`;
}

export default redraw;
