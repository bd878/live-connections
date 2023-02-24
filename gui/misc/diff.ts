import error from '../modules/error';

// has list1 + has list2
function diff(list1: string[], list2: string[]): [string[], string[]] {
  const o1 = new Map(list1.map(v => [v, 1]));
  const o2 = new Map(list2.map(v => [v, 1]));

  for (let i = 0; i < list1.length; i++) {
    const v1 = list1[i];

    if (o1.has(v1) && o2.has(v1)) {
      o1.delete(v1);
      o2.delete(v1);
    }
  }

  return [Array.from(o1.keys()), Array.from(o2.keys())];
}

export default diff;