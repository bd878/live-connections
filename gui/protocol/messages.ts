import log from '../modules/log';
import C from './const';

function makeMouseMoveMessage(x: number, y: number): ABuffer {
  log.Print("x, y:", x, y);

  const messageSize = (
    C.TYPE_SIZE +  // type
    C.COORD_SIZE + // x-coord
    C.COORD_SIZE   // y-coord
  );

  const buffer = new ArrayBuffer(C.SIZE_PREFIX_SIZE + messageSize);
  const dv: any = new DataView(buffer);

  let offset = 0;
  dv.setUint16(offset, messageSize, C.ENDIANNE);
  offset += C.SIZE_PREFIX_SIZE;

  dv.setInt8(offset, C.MOUSE_MOVE_TYPE, C.ENDIANNE);
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
    areaEncoded.size + // area bytes
    C.SIZE_PREFIX_SIZE + // user size
    userEncoded.size   // user bytes
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

export { makeMouseMoveMessage, makeAuthUserMessage };
