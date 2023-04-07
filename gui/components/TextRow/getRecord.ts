import type TextRow from './index';
import error from '../../modules/error';

function getRecord(this: TextRow): TextRecord {
  if (!this.record) throw error.failedToGet("TextRow/getRecord", "record");
  return this.record;
}

export default getRecord;
