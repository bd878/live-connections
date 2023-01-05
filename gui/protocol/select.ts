import C from './const';
import log from 'modules/log';
import {
  parseMouseMoveMessage,
  parseAuthOkMessage,
  parseUsersOnlineMessage,
} from './parser';
import {
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
  onAuthOk,
} from './handlers';

async function select(event: any /* another set of bytes have come... */ ) {
  let buffer = await event.data.arrayBuffer();
  const dv = new DataView(buffer);

  let offset = 0;
  let size = 0;

  while (offset <= size) {
    size = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    if (size === 0) {
      throw new Error('[select]: size is 0 =', size);
    }

    const type = dv.getInt8(offset, C.ENDIANNE);
    offset += C.TYPE_SIZE;

    const slice = buffer.slice(offset);

    switch (type) {
      case C.MOUSE_MOVE_TYPE:
        setTimeout(() => onMouseMove(
          parseMouseMoveMessage(slice)
        ), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.INIT_MOUSE_COORDS_TYPE:
        setTimeout(() => onInitMouseCoords(
          parseMouseMoveMessage(slice)
        ), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.AUTH_OK_TYPE:
        const message = new Blob([slice]);
        setTimeout(() => onAuthOk(
          parseAuthOkMessage(message)
        ), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.USERS_ONLINE_TYPE:
        setTimeout(() => onUsersOnline(
          parseUsersOnlineMessage(slice)
        ), 0); /* throw it in a loop */
        offset += size;
        break;
      default:
        log.Print("[select]: unknown type =", type);
        return;
    }
  }
}

export default select;
