import Log from '../modules/log';
import { isIdentifable } from '../rtti';
import getUid from '../misc/getUid';

const log = new Log("methods/free");

function free(this: Elem): boolean {
  if (isIdentifable(this)) log.Debug(getUid(this.getName(), this.id));
  else log.Warn("trying to free unindentifable elem of", this.getName());

  return true;
}

export default free;
