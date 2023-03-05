import { makeAuthUserMessage } from './messages';
import { log } from '../modules/log';
import socket from '../net/socket';

async function authUser(areaName: AreaName, userName: UserName) {
  log.Debug("authUser", "auth user");

  const authMessage = await makeAuthUserMessage(areaName, userName);
  socket.send(authMessage);
}

async function establish(areaName: AreaName, userName: UserName) {
  log.Debug("establish", "establish");

  await socket.waitOpen();
  await authUser(areaName, userName);

  // run
  log.Debug("establish", "socket is running..."); // DEBUG

  if (!socket.isReady()) {
    log.Debug("establish", "error on message handling");
    throw new Error("socket is not ready");
  }
}

export default establish;
