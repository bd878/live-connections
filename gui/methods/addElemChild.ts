import addElem from './addElem';
import { isAccessible } from '../rtti';
import error from '../modules/error';

function addElemChild(this: Elem & Containable, key: string, elem: Elem): boolean {
  const res = addElem.call(this, key, elem);
  if (!res) {
    throw error.create('addElemChild', 'failed to add elem to container =', key);
  }

  if (!isAccessible(this)) {
    throw error.wrongInterface('addElemChild', "target elem does not implement Accessible");
  }

  if (!isAccessible(elem)) {
    throw error.wrongInterface('addElemChild', "value elem does not implement Accessible");
  }

  this.get().appendChild(elem.get());
  return true;
}

export default addElemChild;
