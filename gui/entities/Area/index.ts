import { create } from './static';

const areas: Area[] = [];

class Area {}

function make(): Area {
  const area = new Area();
  areas.push(area)
  return area;
}

export default {
  create,
  make
};
