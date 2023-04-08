import type Square from './index';
import TextArea from '../TextArea';
import TitlesList from '../TitlesList';
import create from '../../methods/create';

function createSquare(this: Square, id: Id = ''): HTMLElement {
  const titlesList = new TitlesList();
  titlesList.create(id);

  const textarea = new TextArea();
  textarea.create(id);

  const root = create.call(this, id);
  this.addElem(titlesList.getUid(), titlesList);
  this.addElem(textarea.getUid(), textarea);

  return root;
}

export default createSquare;
