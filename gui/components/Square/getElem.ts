import type Square from './index';
import error from '../../modules/error';
import getElem from '../../methods/getElem';
import { isAccessible } from '../../rtti';

function getAccessibleElem(this: Square, uid: Uid): Elem & Accessible {
  const elem = getElem.call(this, uid);
  if (!isAccessible(elem)) throw error.wrongInterface("square.getElem", uid, "not accessible");
  return elem;
}

export default getAccessibleElem;
