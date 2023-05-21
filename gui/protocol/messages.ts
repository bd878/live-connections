import Log from '../modules/log';
import C from './const';

const log = new Log("messages");

function makeMouseMoveMessage(x: number, y: number): ABuffer {
  return makeCoordsMessage(x, y, C.MOUSE_MOVE_TYPE);
}

function makeSquareMoveMessage(x: number, y: number): ABuffer {
  return makeCoordsMessage(x, y, C.SQUARE_MOVE_TYPE);
}

// messageSize + type
async function makeAddRecordMessage(): Promise<ABuffer> {
  const messageSize = (
    C.TYPE_SIZE // type
  );

  const buffer = new ArrayBuffer(C.SIZE_PREFIX_SIZE + messageSize);
  const dv: any = new DataView(buffer);

  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, C.ADD_RECORD_TYPE, C.ENDIANNE);
  offset += C.TYPE_SIZE;

  return buffer;
}

// messageSize + type + textSize + text
async function makeTextInputMessage(text: string): Promise<ABuffer> {
  const textEncoded = new Blob([text], { type: "text/plain"});
  const textBuffer = await textEncoded.arrayBuffer();
  const typedText = new Uint8Array(textBuffer);

  const messageSize = (
    C.TYPE_SIZE        + // type
    C.SIZE_PREFIX_SIZE + // text size
    textEncoded.size     // text
  );

  const buffer = new ArrayBuffer(C.SIZE_PREFIX_SIZE + messageSize);
  const dv: any = new DataView(buffer);

  // message
  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, C.TEXT_INPUT_TYPE, C.ENDIANNE);
  offset += C.TYPE_SIZE;

  // text
  dv.setUint16(offset, textEncoded.size, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  for (let i = 0; i < typedText.length; i++, offset++) {
    dv.setUint8(offset, typedText[i], C.ENDIANNE);
  }

  return buffer;
}

function makeCoordsMessage(x: number, y: number, messageType: number): ABuffer {
  const messageSize = (
    C.TYPE_SIZE  + // type
    C.COORD_SIZE + // x-coord
    C.COORD_SIZE   // y-coord
  );

  const buffer = new ArrayBuffer(C.SIZE_PREFIX_SIZE + messageSize);
  const dv: any = new DataView(buffer);

  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, messageType, C.ENDIANNE);
  offset += C.TYPE_SIZE;

  dv.setFloat32(offset, x, C.ENDIANNE);
  offset += C.COORD_SIZE;

  dv.setFloat32(offset, y, C.ENDIANNE);
  offset += C.COORD_SIZE;

  return buffer;
}

async function makeAuthUserMessage(area: AreaName, user: UserName): Promise<ABuffer> {
  const areaEncoded = new Blob([area], { type: "text/plain"});
  const userEncoded = new Blob([user], { type: "text/plain"});

  const areaArrayBuffer = await areaEncoded.arrayBuffer();
  const userArrayBuffer = await userEncoded.arrayBuffer();

  const typedArea = new Uint8Array(areaArrayBuffer);
  const typedUser = new Uint8Array(userArrayBuffer);

  const messageSize = (
    C.TYPE_SIZE        + // type
    C.SIZE_PREFIX_SIZE + // area size
    areaEncoded.size   + // area bytes
    C.SIZE_PREFIX_SIZE + // user size
    userEncoded.size     // user bytes
  );

  const buffer = new ArrayBuffer(
    C.SIZE_PREFIX_SIZE + // total size
    messageSize
  );
  const dv: any = new DataView(buffer);

  // message
  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, C.AUTH_USER_TYPE, C.ENDIANNE);
  offset += C.TYPE_SIZE;

  // area
  dv.setUint16(offset, areaEncoded.size, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  for (let i = 0; i < typedArea.length; i++, offset++) {
    dv.setUint8(offset, typedArea[i], C.ENDIANNE);
  }

  // user
  dv.setUint16(offset, userEncoded.size, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  for (let i = 0; i < typedUser.length; i++, offset++) {
    dv.setUint8(offset, typedUser[i], C.ENDIANNE);
  }

  return buffer;
}

function makeSelectRecordMessage(recordId: number): ABuffer {
  const messageSize = (
    C.TYPE_SIZE + // type
    C.ID_SIZE     // id size
  );

  const buffer = new ArrayBuffer(
    C.SIZE_PREFIX_SIZE + // total size
    messageSize
  );

  const dv: any = new DataView(buffer);

  // message
  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, C.SELECT_RECORD_TYPE, C.ENDIANNE);
  offset += C.TYPE_SIZE;

  dv.setInt32(offset, recordId, C.ENDIANNE);
  offset += C.COORD_SIZE;

  return buffer;
}

export {
  makeMouseMoveMessage,
  makeAuthUserMessage,
  makeSquareMoveMessage,
  makeTextInputMessage,
  makeAddRecordMessage,
  makeSelectRecordMessage,
};
