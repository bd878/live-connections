import log from "./log.js";
import C from "./constants.js";
import { makeAuthUserMessage } from "./messages.js";

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

async function authUser(socket, user) {
  const authMessage = await makeAuthUserMessage(user.area, user.name);
  socket.send(authMessage);
}

function closeHandler(event) {
  ;(event.wasClean
    ? log.Print(`Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Print("Connection died")
  );
}

function errorHandler(event) {
  log.Print("error =", event);
}

async function handleMouseMoveMessage(fn, buf /* ArrayBuffer */) {
  let offset = 0;
  const dv = new DataView(buf)

  let nameSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (nameSize === 0) {
    throw new Error('[handleMouseMoveMessage]: nameSize is 0 =', nameSize);
  }

  const nameBytes = new Uint8Array(buf, offset, nameSize);
  const blob = new Blob([nameBytes]);
  const name = await blob.text();
  offset += nameSize;

  const xPos = dv.getFloat32(offset,  C.ENDIANNE);
  offset += C.COORD_SIZE;

  const yPos = dv.getFloat32(offset, C.ENDIANNE);
  offset += C.COORD_SIZE;

  fn({ name, xPos, yPos });
}

async function handleAuthOkMessage(fn, message /* Blob */) {
  const text = await message.text();
  fn(text);
}

function handleUsersOnlineMessage(fn, buffer /* ArrayBuffer */) {
  log.Print("handle users online message");
  const users = [];
  fn(users)
}

async function establishProtocol(handlers, socket, user) {
  // establish
  socket.create(C.BACKEND_URL + C.SOCKET_PATH);

  await socket.waitOpen();
  await authUser(socket, user);

  socket.pushSendGuard(user.isNotAuthed)

  // run
  if (socket.isReady()) {
    log.Print("socket is running..."); // DEBUG

    socket.onMessage((event) => messageLoop(handlers, event));
    socket.onClose(closeHandler);
    socket.onError(errorHandler);
  } else {
    throw new Error("[init]: failed to open socket");
  }
}

export {
  establishProtocol,
};
