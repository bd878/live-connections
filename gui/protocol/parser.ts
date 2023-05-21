import Log from '../modules/log';
import C from './const';

const log = new Log("parser");

// userSize + userBytes + recordsCount + []{recordSize + valueSize + valueBytes + titleSize + titleBytes + idBytes + updatedAt + createdAt}
async function parseRecordsListMessage(buf: any /* ArrayBuffer */): Promise<RecordsListEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let userSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (userSize === 0) {
    throw new Error(`[parseRecordsListMessage]: userSize is 0`);
  }

  const userBytes = new Uint8Array(buf, offset, userSize);
  const blob = new Blob([userBytes]);
  const name = await blob.text();
  offset += userSize;

  let recordsCount = dv.getUint16(offset, C.ENDIANNE)
  offset += C.COUNT_USERS_SIZE;

  const records: TextRecord[] = [];
  for (let i = 0; i < recordsCount; i++) {
    const recordSize = dv.getUint16(offset, C.ENDIANNE); // not used
    offset += C.SIZE_PREFIX_SIZE;

    const valueSize = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    const valueBytes = new Uint8Array(buf, offset, valueSize);
    const valueBlob = new Blob([valueBytes]);
    const value = await valueBlob.text();
    offset += valueSize;

    const titleSize = dv.getUint16(offset, C.ENDIANNE);
    offset += C.SIZE_PREFIX_SIZE;

    const titleBytes = new Uint8Array(buf, offset, titleSize);
    const titleBlob = new Blob([titleBytes]);
    const title = await titleBlob.text();
    offset += titleSize;

    const id = dv.getInt32(offset, C.ENDIANNE);
    offset += C.ID_SIZE;

    const updatedAt = dv.getInt32(offset, C.ENDIANNE);
    offset += C.TIMESTAMP_SIZE;

    const createdAt = dv.getInt32(offset, C.ENDIANNE);
    offset += C.TIMESTAMP_SIZE;

    log.Debug("value, title, id, updatedAt, createdAt:", value, title, id, updatedAt, createdAt);

    records.push({ value, title, id, updatedAt, createdAt });
  }

  return { name, records } as RecordsListEvent;
}

// userSize + userBytes + idBytes
async function parseSelectRecordMessage(buf: any /* ArrayBuffer */): Promise<SelectRecordEvent> {
  let offset = 0;
  const dv = new DataView(buf);

  let userSize = dv.getUint16(offset, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  if (userSize === 0) {
    throw new Error(`[parseSelectRecordMessage]: userSize is 0`);
  }

  const userBytes = new Uint8Array(buf, offset, userSize);
  const blob = new Blob([userBytes]);
  const name = await blob.text();
  offset += userSize;

  const id = dv.getInt32(offset, C.ENDIANNE);
  offset += C.ID_SIZE;

  return { name, id } as SelectRecordEvent;
}

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
  parseRecordsListMessage,
  parseSelectRecordMessage,
};
