import C from './const';
import Log from '../modules/log';
import {
  parseCoordsMessage,
  parseAuthOkMessage,
  parseTextInputMessage,
  parseUsersOnlineMessage,
} from './parser';
import {
  onMouseMove,
  onSquareMove,
  onInitSquareCoords,
  onUsersOnline,
  onAuthOk,
  onTextInput,
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
      case C.MOUSE_MOVE_TYPE:
        log.Debug("mouse move");

        setTimeout(() => {
          parseCoordsMessage(slice).then(onMouseMove);
        }, 0); /* throw it in a loop */
        offset += size;
        break;
      case C.SQUARE_MOVE_TYPE:
        log.Debug("square move");

        setTimeout(() => {
          parseCoordsMessage(slice).then(onSquareMove);
        }, 0); /* throw it in a loop */
        offset += size;
        break;
      case C.TEXT_INPUT_TYPE:
        log.Debug("text input");

        setTimeout(() => {
          parseTextInputMessage(slice).then(onTextInput);
        });
      case C.INIT_SQUARE_COORDS_TYPE:
        log.Debug("init square coords");

        setTimeout(() => {
          parseCoordsMessage(slice).then(onInitSquareCoords);
        }, 0); /* throw it in a loop */
        offset += size;
        break;
      case C.AUTH_OK_TYPE:
        log.Debug("auth ok");

        const message = new Blob([slice]);
        setTimeout(() => {
          parseAuthOkMessage(message).then(onAuthOk);
        }, 0); /* throw it in a loop */
        offset += size;
        break;
      case C.USERS_ONLINE_TYPE:
        log.Debug("users online");

        setTimeout(() => {
          parseUsersOnlineMessage(slice).then(onUsersOnline);
        }, 0); /* throw it in a loop */
        offset += size;
        break;
      default:
        log.Debug("unknown type =", type);
        return;
    }
  }
}

export default select;
