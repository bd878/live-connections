import type TextArea from './index';
import error from '../../modules/error';

function redraw(this: TextArea, ...args: any[]) {
  const text = args[0];
  if (!this.textarea) throw error.noElementCreated("textarea/redraw", "no textarea");
  this.textarea = text;
}

export default redraw;
