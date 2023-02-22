import { create } from './static';
import error from '../../modules/error';

let area: Area | null = null;

class Area {
  constructor(public name: string = '') {}
}

function make(areaName: string): Area {
  area = new Area(areaName);
  return area;
}

function getMy(): Area {
  if (!area) {
    throw error.noElementCreated("areas");
  }
  return area;
}

export default {
  create,
  make,
  getMy,
};
