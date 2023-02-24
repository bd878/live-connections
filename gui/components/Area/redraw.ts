import type { Area } from './index';
import rtti from '../../rtti';
import cursors from '../../entities/cursors';
import error from '../../modules/error';

function redraw(this: Area, piece: string, ...args: any[]) {
  if (piece === 'cursor' && args[0]) {
    redrawSingleCursor.call(this, args[0]);
  }
}

function redrawSingleCursor(this: Area, name: CursorName) {
  if (!this.hasElem(name)) {
    throw error.noElementCreated("Area redrawSingleCursor");
  }
  const elem = this.getElem(name);
  if (!rtti.isMovable(elem)) {
    throw error.wrongInterface("Area redrawSingleCursor", name, "is not movable");
  }
  const x = cursors.getX(name);
  const y = cursors.getY(name);
  elem.move(x, y);
}

export default redraw;
