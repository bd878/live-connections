import type { Area } from './index';
import rtti from '../../rtti';
import log from '../../modules/log';
import coords from '../../entities/coords';
import error from '../../modules/error';

function redraw(this: Area, piece: string, ...args: any[]) {
  if (piece === 'cursor' && args[0]) {
    redrawCoords.call(this, args[0]);
  } else if (piece === 'square' && args[0]) {
    redrawCoords.call(this, args[0]);
  }
}

function redrawCoords(this: Area, uid: Uid) {
  log.Print("redrawCoords", uid);
  if (!this.hasElem(uid)) {
    throw error.noElementCreated("Area redrawCoords");
  }
  const elem = this.getElem(uid);
  if (!rtti.isMovable(elem)) {
    throw error.wrongInterface("Area redrawCoords", uid, "is not movable");
  }
  const x = coords.getX(uid);
  const y = coords.getY(uid);
  elem.move(x, y);
}

export default redraw;
