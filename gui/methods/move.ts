import error from '../modules/error';
import { log } from '../modules/log';

/* public */
function move(this: Elem & Moveable, x: number, y: number): void {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  log.Debug("move", `name, x, y: ${this.getName()}, ${x}, ${y}`);

  this.root.style.transform = `
    translate3D(
      ${x}px,
      ${y}px,
      0
    )
  `;
}

export default move;