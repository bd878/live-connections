import C from './const';
import log from '../modules/log';

async function parseMouseMoveMessage(buf: any /* ArrayBuffer */): Promise<MouseMoveEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let nameSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (nameSize === 0) {
    throw new Error(`[parseMouseMoveMessage]: nameSize is 0`);
  }

  const nameBytes = new Uint8Array(buf, offset, nameSize);
  const blob = new Blob([nameBytes]);
  const name = await blob.text();
  offset += nameSize;

  const xPos = dv.getFloat32(offset,  C.ENDIANNE);
  offset += C.COORD_SIZE;

  const yPos = dv.getFloat32(offset, C.ENDIANNE);
  offset += C.COORD_SIZE;

  return { name, xPos, yPos } as MouseMoveEvent;
}

async function parseAuthOkMessage(message: any /* Blob */): Promise<AuthOkEvent> {
  const text: string = await message.text();
  return {text};
}

async function parseUsersOnlineMessage(buffer: any /* ArrayBuffer */): Promise<UsersOnlineEvent> {
  log.Print("parseUsersOnlineMessage", "handle users online message");
  const users: UserName[] = [];
  const colors: string[] = [];
  return {users, colors};
}

export {
  parseMouseMoveMessage,
  parseAuthOkMessage,
  parseUsersOnlineMessage,
};
