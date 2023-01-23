import type { Area } from './index';
import create from '../../methods/create';
import addCSSClass from '../../methods/addCSSClass';

function createArea(this: Area): HTMLElement {
  const result = create.call(this);

  addCSSClass.call(this, 'area');

  return result;
}

export default createArea;
