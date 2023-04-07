import type Square from './index';
import TextArea from '../TextArea';
import TitlesList from '../TitlesList';
import create from '../../methods/create';

function createSquare(this: Square, id: Id = ''): HTMLElement {
  const textarea = new TextArea();
  textarea.create(id);

  const titlesList = new TitlesList();
  titlesList.create(id);

  const root = create.call(this, id);
  this.addElem(textarea.getUid(), textarea);
  this.addElem(titlesList.getUid(), titlesList);

  return root;
}

export default createSquare;
