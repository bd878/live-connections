import { makeAuthUserMessage } from './messages';
import Log from '../modules/log';
import socket from '../net/socket';

const log = new Log("gui/protocol");

async function authUser(areaName: AreaName, userName: UserName) {
  const authMessage = await makeAuthUserMessage(areaName, userName);
  socket.send(authMessage);
}

async function establish(areaName: AreaName, userName: UserName) {
  log.Debug("establish");

  await socket.waitOpen();
  await authUser(areaName, userName);

  // run
  log.Info("socket is running..."); // DEBUG

  if (!socket.isReady()) {
    log.Fail("error on message handling");
    throw new Error("socket is not ready");
  }
}

export default establish;
