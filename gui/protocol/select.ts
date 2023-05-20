import C from './const';
import Log from '../modules/log';
import {
  parseCoordsMessage,
  parseAuthOkMessage,
  parseTextInputMessage,
  parseUsersOnlineMessage,
  parseRecordsListMessage,
} from './parser';
import {
  onMouseMove,
  onSquareMove,
  onInitSquareCoords,
  onUsersOnline,
  onAuthOk,
  onTextInput,
  onListRecords,
} from './handlers';

const log = new Log("protocol/select");

async function select(b: any /* another set of bytes have come... */ ) {
  let buffer = await b.data.arrayBuffer();
  const dv: any = new DataView(buffer);

  let offset = 0;
  let size = 0;

  while (offset <= size) {
    size = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    if (size === 0) {
      throw new Error(`[select]: size is 0`);
    }

    const type = dv.getInt8(offset, C.ENDIANNE);
    offset += C.TYPE_SIZE;

    const slice = buffer.slice(offset);

    switch (type) {
      case C.RECORDS_LIST_TYPE:
        setTimeout(() => {
          parseRecordsListMessage(slice).then(onListRecords);
        }, 0); /* throw it in a loop */
        break;
      case C.MOUSE_MOVE_TYPE:
        setTimeout(() => {
          parseCoordsMessage(slice).then(onMouseMove);
        }, 0); /* throw it in a loop */
        break;
      case C.SQUARE_MOVE_TYPE:
        setTimeout(() => {
          parseCoordsMessage(slice).then(onSquareMove);
        }, 0); /* throw it in a loop */
        break;
      case C.TEXT_INPUT_TYPE:
        setTimeout(() => {
          parseTextInputMessage(slice).then(onTextInput);
        }, 0); /* throw it in a loop */
        break;
      case C.INIT_SQUARE_COORDS_TYPE:
        setTimeout(() => {
          parseCoordsMessage(slice).then(onInitSquareCoords);
        }, 0); /* throw it in a loop */
        break;
      case C.AUTH_OK_TYPE:
        setTimeout(() => {
          parseAuthOkMessage(new Blob([slice])).then(onAuthOk);
        }, 0); /* throw it in a loop */
        break;
      case C.USERS_ONLINE_TYPE:
        setTimeout(() => {
          parseUsersOnlineMessage(slice).then(onUsersOnline);
        }, 0); /* throw it in a loop */
        break;
      default:
        log.Debug("unknown type =", type);
        return;
    }

    offset += size;
  }
}

export default select;
