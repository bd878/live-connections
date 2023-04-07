import error from '../modules/error';

const defaultMapper = (v: string) => v;

// has list1 + has list2
function diff<T = string, R = string>(
  list1: T[],
  list2: T[],
  mapper: ((v: T) => R),
): [T[], T[]] {
  const o1 = new Map(list1.map(v => [mapper(v), v]));
  const o2 = new Map(list2.map(v => [mapper(v), v]));

  for (let i = 0; i < list1.length; i++) {
    const v1 = list1[i];
    const k1 = mapper(v1);

    if (o1.has(k1) && o2.has(k1)) {
      o1.delete(k1);
      o2.delete(k1);
    }
  }

  return [Array.from(o1.values()), Array.from(o2.values())];
}

export default diff;
export { defaultMapper };
