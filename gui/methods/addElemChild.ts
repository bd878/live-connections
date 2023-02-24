import addElem from './addElem';
import rtti from '../rtti';
import error from '../modules/error';

function addElemChild(this: Elem & Containable, key: string, elem: Elem): boolean {
  const res = addElem.call(this, key, elem);
  if (!res) {
    throw error.create('addElemChild', 'failed to add elem to container =', key);
  }

  if (!rtti.isAccessible(this)) {
    throw error.wrongInterface('addElemChild', "target elem does not implement Accessible");
  }

  if (!rtti.isAccessible(elem)) {
    throw error.wrongInterface('addElemChild', "value elem does not implement Accessible");
  }

  this.get().appendChild(elem.get());
  return true;
}

export default addElemChild;
