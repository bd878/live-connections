import type { UsersPanel } from './index';
import UsersList from '../UsersList';
import create from '../../methods/create';
import append from '../../methods/append';

function createUsersPanel(this: UsersPanel): HTMLElement {
  UsersList.create();

  const result = create.call(this);
  append.call(this, UsersList);

  return result;
}

export default createUsersPanel;
