import type TitlesList from './index';
import TextRow from '../TextRow';
import Log from '../../modules/log';

const log = new Log("TitlesList/listRecords");

function listRecords(this: TitlesList): TextRecord[] {
  const result: TextRecord[] = [];
  for (const value of this.container.values()) {
    if (value instanceof TextRow) {
      result.push(value.getRecord());
    } else {
      log.Warn("not a TextRow instance");
    }
  }
  return result;
}

export default listRecords;
