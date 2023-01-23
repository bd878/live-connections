import type { Main } from './index';
import Area from '../Area';
import UsersPanel from '../UsersPanel';
import create from '../../methods/create';
import append from '../../methods/append';
import addCSSClass from '../../methods/addCSSClass';

function createMain(this: Main): HTMLElement {
  UsersPanel.create();
  Area.create();

  const result = create.call(this);
  addCSSClass.call(this, 'main');
  append.call(this, UsersPanel);
  append.call(this, Area);

  return result;
}

export default createMain;
