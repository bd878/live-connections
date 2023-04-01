import type TextArea from './index';
import error from '../../modules/error';

function turnReadonly(this: TextArea) {
  if (!this.root) throw error.noElementCreated("textarea/turnReadonly", "no root");
  this.root.setAttribute("readonly", "1");
}

export default turnReadonly;
