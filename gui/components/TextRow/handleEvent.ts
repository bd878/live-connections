import type TextRow from './index';
import Log from '../../modules/log';
import { onSelectRecord } from '../../listeners';

const log = new Log("TextRow/handleEvent");

function handleEvent(this: TextRow, event: any) {
  if (event.type === "click") {
    onSelectRecord(this.getRecord().id);
  } else {
    log.Warn("no handler for event", event.type);
  }
}

export default handleEvent;
