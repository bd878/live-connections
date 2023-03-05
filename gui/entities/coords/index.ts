import error from '../../modules/error';

let coordsX: Map<Uid, number> = new Map();
let coordsY: Map<Uid, number> = new Map();
let uids: Uid[] = [];
let total: number = 0;

function length(): number {
  return total;
}

function getUidAt(i: number): Uid {
  const uid = uids[i];
  if (!uid) {
    throw error.outOfRange("coords", i);
  }
  return uid;
}

function has(uid: Uid): boolean {
  return coordsX.has(uid);
}

function remove(uid: Uid) {
  coordsX.delete(uid);
  coordsY.delete(uid);
  uids = uids.filter(n => n != uid);
  total--;
}

function set(uid: Uid, xPos: number, yPos: number) {
  if (!coordsX.has(uid) && !coordsY.has(uid)) {
    uids.push(uid);
    total++;
  }

  coordsX.set(uid, xPos);
  coordsY.set(uid, yPos);
}

function getX(uid: Uid): number {
  const xPos = coordsX.get(uid);
  if (xPos === undefined) throw error.failedToGet("coords getX", uid);
  return xPos;
}

function getY(uid: Uid): number {
  const yPos = coordsY.get(uid);
  if (yPos === undefined) throw error.failedToGet("coords getY", uid);
  return yPos;
}

export default {
  length,
  getUidAt,
  has,
  set,
  getX,
  getY,
  remove,
};
