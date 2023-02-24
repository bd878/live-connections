import type { Area } from './index';
import rtti from '../../rtti';
import cursors from '../../entities/cursors';
import error from '../../modules/error';
import Cursor from '../Cursor';

function redraw(this: Area, piece: string, ...args: any[]) {
  if (piece === 'cursors' || piece === '') {
    redrawCursors.call(this);
  } else if (piece === 'cursor' && args[0]) {
    redrawSingleCursor.call(this, args[0]);
  }
}

function redrawCursors(this: Area) {
  for (let i = 0; i < cursors.length(); i++) {
    const name = cursors.getNameAt(i);
    if (!this.hasElem(name)) {
      const cursor = new Cursor();
      cursor.create();
      this.addElem(name, cursor);
    }

    redrawSingleCursor.call(this, name);
  }
}

function redrawSingleCursor(this: Area, name: UserName) {
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
