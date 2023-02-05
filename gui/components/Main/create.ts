import type { Main } from './index';
import Area from '../Area';
import UsersPanel from '../UsersPanel';
import create from '../../methods/create';
import append from '../../methods/append';

function createMain(this: Main): HTMLElement {
  UsersPanel.create();
  Area.create();

  const result = create.call(this);
  append.call(this, UsersPanel);
  append.call(this, Area);

  return result;
}

export default createMain;
