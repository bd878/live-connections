import type TextRow from './index';
import error from '../../modules/error';

function redraw(this: TextRow, ...args: any[]) {
  if (!this.root) throw error.noElementCreated("textrow/redraw", "no root");
  if (!this.record) throw error.failedToGet("textrow/redraw", "no record");
  this.root.textContent = `${this.record.id}`;
}

export default redraw;
