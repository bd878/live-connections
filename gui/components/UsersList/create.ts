import type { UsersList } from './index';
import create from '../../methods/create';
import addCSSClass from '../../methods/addCSSClass';

function createUsersList(this: UsersList): HTMLElement {
  const result = create.call(this);

  addCSSClass.call(this, 'users-list');

  return result;
}

export default createUsersList;
