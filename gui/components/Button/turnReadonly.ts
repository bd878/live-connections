import type Button from './index';
import error from '../../modules/error';

function turnReadonly(this: Button) {
  if (!this.root) throw error.noElementCreated("button/turnReadonly", "no root");
  this.root.setAttribute("disabled", "1");
}

export default turnReadonly;
