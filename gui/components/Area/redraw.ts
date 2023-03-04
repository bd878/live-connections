import type { Area } from './index';
import rtti from '../../rtti';
import log from '../../modules/log';
import coords from '../../entities/coords';
import error from '../../modules/error';

function redraw(this: Area, piece: string, ...args: any[]) {
  if (piece === 'cursor' && args[0]) {
    redrawSingleCursor.call(this, args[0]);
  }
}

function redrawSingleCursor(this: Area, uid: Uid) {
  log.Print("redrawSingleCursor", uid);
  if (!this.hasElem(uid)) {
    throw error.noElementCreated("Area redrawSingleCursor");
  }
  const elem = this.getElem(uid);
  if (!rtti.isMovable(elem)) {
    throw error.wrongInterface("Area redrawSingleCursor", uid, "is not movable");
  }
  const x = coords.getX(uid);
  const y = coords.getY(uid);
  elem.move(x, y);
}

export default redraw;
