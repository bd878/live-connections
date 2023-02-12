import type { Area } from './index';
import addElem from '../../methods/addElem';
import error from '../../modules/error';

function addElemOnArea(this: Area, key: string, elem: Elem & Accessible): boolean {
  const res = addElem.call(this, key, elem);
  if (!res) {
    throw error.create('addElemOnArea', 'failed to add elem to container =', key);
  }

  this.get().appendChild(elem.get());
  return true;
}

export default addElemOnArea;
