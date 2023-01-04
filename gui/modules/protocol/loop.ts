
async function messageLoop(hs, event /* another set of bytes have come... */ ) {
  let buffer = await event.data.arrayBuffer();
  const dv = new DataView(buffer);

  let offset = 0;
  let size = 0;

  while (offset <= size) {
    size = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    if (size === 0) {
      throw new Error('[messageLoop]: size is 0 =', size);
    }

    const type = dv.getInt8(offset, C.ENDIANNE);
    offset += C.TYPE_SIZE;

    const slice = buffer.slice(offset);

    switch (type) {
      case C.MOUSE_MOVE_TYPE:
        setTimeout(() => handleMouseMoveMessage(hs.onMouseMove, slice), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.INIT_MOUSE_COORDS_TYPE:
        setTimeout(() => handleMouseMoveMessage(hs.onInitMouseCoords, slice), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.AUTH_OK_TYPE:
        const message = new Blob([slice]);
        setTimeout(() => handleAuthOkMessage(hs.onAuthOk, message), 0); /* throw it in a loop */
        offset += size;
        break;
      case C.USERS_ONLINE_TYPE:
        setTimeout(() => handleUsersOnlineMessage(hs.onUsersOnline, slice), 0); /* throw it in a loop */
        offset += size;
        break;
      default:
        log.Print("[messageLoop]: unknown type =", type);
        return;
    }
  }
}

export default messageLoop;
