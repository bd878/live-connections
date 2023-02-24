import error from '../modules/error';
import rtti from '../rtti';
import delElem from './delElem';

function delElemChild(this: Elem & Containable, key: string) {
  const elem = this.getElem(key);

  delElem.call(this, key);

  if (!rtti.isAccessible(this)) {
    throw error.wrongInterface('addElemChild', "target elem does not implement Accessible");
  }

  if (!rtti.isAccessible(elem)) {
    throw error.wrongInterface('addElemChild', "target elem does not implement Accessible");
  }

  this.get().removeChild(elem.get());
}

export default delElemChild;
