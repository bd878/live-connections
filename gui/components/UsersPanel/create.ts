import type { UsersPanel } from './index';
import UsersList from '../UsersList';
import create from '../../methods/create';
import append from '../../methods/append';

function createUsersPanel(this: UsersPanel, id: Id = ''): HTMLElement {
  UsersList.create(id);

  const result = create.call(this, id);
  append.call(this, UsersList);

  return result;
}

export default createUsersPanel;
