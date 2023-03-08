import type TextArea from './index';
import create from '../../methods/create';

function createTextArea(this: TextArea): HTMLElement {
  const root = create.call(this);
  this.textarea = document.createElement("textarea");
  return root;
}

export default createTextArea;
