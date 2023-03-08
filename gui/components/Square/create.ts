import type Square from './index';
import TextArea from '../TextArea';
import create from '../../methods/create';

function createSquare(this: Square, id: Id = ''): HTMLElement {
  const textarea = new TextArea();
  textarea.create(id);

  const root = create.call(this, id);
  this.addElem(textarea.getUid(), textarea);

  return root;
}

export default createSquare;
