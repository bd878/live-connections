import { create } from './static';
import error from '../../modules/error';

let _myName: AreaName | null = null;

function setMyName(areaName: string) {
  _myName = areaName;
}

function myName(): AreaName {
  if (!_myName) throw error.failedToGet("areas myName");
  return _myName;
}

export default {
  create,
  setMyName,
  myName,
};
