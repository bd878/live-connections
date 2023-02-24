import error from '../../modules/error';

let cursorsX: Map<CursorName, number> = new Map();
let cursorsY: Map<CursorName, number> = new Map();
let names: CursorName[] = []; // TODO: nameToIndex, indextToName
let total: number = 0;

function length(): number {
  return total;
}

function getNameAt(i: number): CursorName {
  const name = names[i];
  if (!name) {
    throw error.outOfRange("cursors", i);
  }
  return name;
}

function has(name: CursorName): boolean {
  return cursorsX.has(name);
}

function remove(name: CursorName) {
  cursorsX.delete(name);
  cursorsY.delete(name);
  names = names.filter(n => n != name);
  total--;
}

function set(name: CursorName, xPos: number, yPos: number) {
  if (!cursorsX.has(name) && !cursorsY.has(name)) {
    names.push(name);
    total++;
  }

  cursorsX.set(name, xPos);
  cursorsY.set(name, yPos);
}

function getX(name: CursorName): number {
  const xPos = cursorsX.get(name);
  if (!xPos) throw error.failedToGet("cursors getX", name);
  return xPos;
}

function getY(name: CursorName): number {
  const yPos = cursorsY.get(name);
  if (!yPos) throw error.failedToGet("cursors getY", name);
  return yPos;
}

export default {
  length,
  getNameAt,
  has,
  set,
  getX,
  getY,
  remove,
};
