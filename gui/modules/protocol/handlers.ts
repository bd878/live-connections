import C from './constants';
import log from 'modules/log';

async function handleMouseMoveMessage(fn: Fn, buf: any /* ArrayBuffer */) {
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

export {
  handleMouseMoveMessage,
  handleAuthOkMessage,
  handleUsersOnlineMessage,
};
