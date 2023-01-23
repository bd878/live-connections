import type { UsersPanel } from './index';
import UsersList from '../UsersList';
import create from '../../methods/create';
import append from '../../methods/append';
import addCSSClass from '../../methods/addCSSClass';

function createUsersPanel(this: UsersPanel): HTMLElement {
  UsersList.create();

  const result = create.call(this);
  addCSSClass.call(this, 'users-panel');
  append.call(this, UsersList);

  return result;
}

export default createUsersPanel;
