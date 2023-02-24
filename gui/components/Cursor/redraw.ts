import type Cursor from './index';
import error from '../../modules/error';

function redraw(this: Cursor, piece: string = '') {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }

  this.root.style.backgroundColor = this.color;
}

export default redraw;
