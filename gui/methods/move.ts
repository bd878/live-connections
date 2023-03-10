import error from '../modules/error';

/* public */
function move(this: Elem & Moveable, x: number, y: number): void {
  if (!this.root) {
    throw error.noElementCreated(this.getName());
  }

  this.root.style.transform = `
    translate3D(
      ${x}px,
      ${y}px,
      0
    )
  `;
}

export default move;