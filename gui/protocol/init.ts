import { makeAuthUserMessage } from './messages';
import log from 'modules/log';
import socket from 'net/socket';

async function authUser(areaName: AreaName, userName: UserName) {
  const authMessage = await makeAuthUserMessage(areaName, userName);
  socket.send(authMessage);
}

async function establish(areaName: AreaName, userName: UserName): void {
  await socket.waitOpen();
  await authUser(areaName, userName);

  // run
  log.Print("socket is running..."); // DEBUG

  if (!socket.isReady()) {
    log.Print("[protocol establish]: error on message handling");
    throw new Error(err);    
  }
}

export default establish;
