import rtti from '../rtti';
import error from '../modules/error';

function addChild(this: Elem & Containable, elem: Elem & Accessible): Elem & Containable {
  if (!rtti.isAccessible(this)) {
    throw error.wrongInterface('addElemChild', "target elem does not implement Accessible");
  }

  if (!rtti.isAccessible(elem)) {
    throw error.wrongInterface('addElemChild', "value elem does not implement Accessible");
  }

  this.get().appendChild(elem.get());
  return this;
}

export default addChild;
