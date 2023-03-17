import Log from '../modules/log';
import C from './const';

const log = new Log("parser");

async function parseCoordsMessage(buf: any /* ArrayBuffer */): Promise<CoordsEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let nameSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (nameSize === 0) {
    throw new Error(`[parseCoordsMessage]: nameSize is 0`);
  }

  const nameBytes = new Uint8Array(buf, offset, nameSize);
  const blob = new Blob([nameBytes]);
  const name = await blob.text();
  offset += nameSize;

  const xPos = dv.getFloat32(offset,  C.ENDIANNE);
  offset += C.COORD_SIZE;

  const yPos = dv.getFloat32(offset, C.ENDIANNE);
  offset += C.COORD_SIZE;

  return { name, xPos, yPos } as CoordsEvent;
}

async function parseAuthOkMessage(message: any /* Blob */): Promise<AuthOkEvent> {
  const text: string = await message.text();
  return { text };
}

// userSize + userBytes + textSize + textBytes
async function parseTextInputMessage(buf: any /* ArrayBuffer */): Promise<TextInputEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let nameSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (nameSize === 0) {
    throw new Error(`[parseTextInputMessage]: nameSize is 0`);
  }

  const nameBytes = new Uint8Array(buf, offset, nameSize);
  const nameBlob = new Blob([nameBytes]);
  const name = await nameBlob.text();
  offset += nameSize;

  let textSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  let text = '';
  if (textSize > 0) {
    const textBytes = new Uint8Array(buf, offset, textSize);
    const textBlob = new Blob([textBytes]);
    text = await textBlob.text();
    offset += textSize;
  }

  return { text, name };
}

async function parseUsersOnlineMessage(buf: any /* ArrayBuffer */): Promise<UsersOnlineEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let usersCount = dv.getUint16(offset, C.ENDIANNE)
  offset += C.COUNT_USERS_SIZE;

  const users: UserName[] = [];
  for (let i = 0; i < usersCount; i++) {
    const nameSize = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    const nameBytes = new Uint8Array(buf, offset, nameSize);
    const blob = new Blob([nameBytes]);
    const name = await blob.text();
    offset += nameSize;

    users.push(name);
  }

  return { users };
}

export {
  parseCoordsMessage,
  parseAuthOkMessage,
  parseUsersOnlineMessage,
  parseTextInputMessage,
};
