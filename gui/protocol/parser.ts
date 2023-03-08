import C from './const';
import Log from '../modules/log';

const log = new Log("protocol/parser");

async function parseCoordsMessage(buf: any /* ArrayBuffer */): Promise<CoordsEvent> {
  log.Debug("parseCoordsMessage");

  let offset = 0;
  const dv = new DataView(buf);

  let nameSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (nameSize === 0) {
    log.Fail("nameSize is 0");
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

async function parseTextInputMessage(message: any /* Blob */): Promise<TextInputEvent> {
  const text: string = await message.text();
  return { text, name: '' };
}

async function parseUsersOnlineMessage(buf: any /* ArrayBuffer */): Promise<UsersOnlineEvent> {
  log.Debug("parseUsersOnlineMessage");

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

  log.Debug("parseUsersOnlineMessage");

  return { users };
}

export {
  parseCoordsMessage,
  parseAuthOkMessage,
  parseUsersOnlineMessage,
  parseTextInputMessage,
};
