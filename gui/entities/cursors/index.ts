import error from '../../modules/error';

let cursorsX: Map<Uid, number> = new Map();
let cursorsY: Map<Uid, number> = new Map();
let uids: Uid[] = [];
let total: number = 0;

function length(): number {
  return total;
}

function getUidAt(i: number): Uid {
  const uid = uids[i];
  if (!uid) {
    throw error.outOfRange("cursors", i);
  }
  return uid;
}

function has(uid: Uid): boolean {
  return cursorsX.has(uid);
}

function remove(uid: Uid) {
  cursorsX.delete(uid);
  cursorsY.delete(uid);
  uids = uids.filter(n => n != uid);
  total--;
}

function set(uid: Uid, xPos: number, yPos: number) {
  if (!cursorsX.has(uid) && !cursorsY.has(uid)) {
    uids.push(uid);
    total++;
  }

  cursorsX.set(uid, xPos);
  cursorsY.set(uid, yPos);
}

function getX(uid: Uid): number {
  const xPos = cursorsX.get(uid);
  if (!xPos) throw error.failedToGet("cursors getX", uid);
  return xPos;
}

function getY(uid: Uid): number {
  const yPos = cursorsY.get(uid);
  if (!yPos) throw error.failedToGet("cursors getY", uid);
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
