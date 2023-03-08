import type { Area } from './index';
import error from '../../modules/error';
import getElem from '../../methods/getElem';
import { isMovable } from '../../rtti';

function getMovableElem(this: Area, uid: Uid): Elem & Moveable {
  const elem = getElem.call(this, uid);
  if (!isMovable(elem)) throw error.wrongInterface("area.getElem", uid, "not movable");
  return elem;
}

export default getMovableElem;
