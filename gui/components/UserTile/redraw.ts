import type UserTile from './index';
import error from '../../modules/error';

function redraw(this: UserTile, piece: string = '') {
  if (!this.root) {
    throw error.noElementCreated(this.name);
  }
  this.root.style.backgroundColor = this.color;
}

export default redraw;