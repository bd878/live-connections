import type Square from './index';
import Button from '../Button';
import TextArea from '../TextArea';
import Container from '../Container';
import TitlesList from '../TitlesList';
import create from '../../methods/create';

function createSquare(this: Square, id: Id = ''): HTMLElement {
  const button = new Button();
  button.create(id);

  const titlesList = new TitlesList();
  titlesList.create(id);

  const container = new Container("square-container");
  container.create();

  const textarea = new TextArea();
  textarea.create(id);

  container
    .addChild(button)
    .addChild(titlesList);

  const root = create.call(this, id);

  this
    .addChild(container)
    .addElem(button.getUid(), button)
    .addElem(titlesList.getUid(), titlesList);

  this
    .addChild(textarea)
    .addElem(textarea.getUid(), textarea);

  return root;
}

export default createSquare;
