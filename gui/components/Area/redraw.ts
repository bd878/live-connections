import type { Area } from './index';
import rtti from '../../rtti';
import { log } from '../../modules/log';
import coords from '../../entities/coords';
import error from '../../modules/error';

function redraw(this: Area, ...args: any[]) {
  const piece = args[0];
  const uid = args[1];
  if (piece === 'cursor' && uid) {
    redrawCoords.call(this, uid);
  } else if (piece === 'square' && uid) {
    redrawCoords.call(this, uid);
  }
}

function redrawCoords(this: Area, uid: Uid) {
  log.Debug("redrawCoords", uid);
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
