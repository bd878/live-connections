import type { Area } from './index';
import areas from '../../entities/areas';
import free from '../../methods/free';

function freeArea(this: Area): boolean {
  const result = free.call(this);
  areas.removeMe();
  return result;
}

export default freeArea;
