import getUid from '../misc/getUid';

function _getUid(this: Elem & Identifable): Uid {
  return getUid(this.getName(), this.id);
}

export default _getUid;
