import type Square from './index';
import TextArea from '../TextArea';
import Container from '../Container';
import TitlesList from '../TitlesList';
import Button from '../Button';
import create from '../../methods/create';

function createSquare(this: Square, id: Id = ''): HTMLElement {
  const titlesList = new TitlesList();
  titlesList.create(id);

  const container = new Container("square-container");
  container.create();

  const textarea = new TextArea();
  textarea.create(id);

  const button = new Button();
  button.create(id);

  container
    .addChild(button)
    .addChild(titlesList);

  const root = create.call(this, id);

  this.addElem(button.getUid(), button)

  this
    .addChild(container)
    .addElem(titlesList.getUid(), titlesList);

  this
    .addChild(textarea)
    .addElem(textarea.getUid(), textarea);

  return root;
}

export default createSquare;
