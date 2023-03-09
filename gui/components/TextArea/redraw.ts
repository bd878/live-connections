import type TextArea from './index';
import error from '../../modules/error';

function redraw(this: TextArea, ...args: any[]) {
  const text = args[0];
  if (!this.root) throw error.noElementCreated("textarea/redraw", "no root");
  this.root.textContent = text;
}

export default redraw;
