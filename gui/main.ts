import { makeMouseMoveMessage } from 'protocol/messages';
import select from 'protocol/select';
import establish from 'protocol/init';
import takeAreaName from 'misc/takeAreaName';
import findUserName from 'misc/findUserName';
import log from 'modules/log';
import User from 'entities/User';
import Area from 'entities/Area';
import socket from 'net/socket';

/* Waits for protocol message on socket */
async function run(): void {
  let resolve, reject;
  const p = new Promise((r, j) => {
    resolve = r;
    reject = j;
  });

  try {
    while (1) {
      const message = await socket.waitMessage();
      select(message);
    }

    ;(resolve && resolve(true));
  } catch (e) {
    log.Print("[main run]: failed to run");
    ;(reject && reject(e));
  }
}

/* Applies to server for new area allocation */
async function proceedNewArea(): AreaName {
  const areaName = await Area.create();
  setUrl(`/${areaName}`);
  return areaName;
}

/* Applies to server for new user registration */
async function proceedNewUser(areaName: AreaName): UserName {
  const userName = await User.create(areaName);
  bindUserToArea(areaName, userName);
  return userName;
}

/* Initializes internal parts: area, user, socket, protocol etc. */
async function main() {
  log.mode = 'debug';
  socket.init();

  let userName;
  let areaName = takeAreaName(window.location.pathname);

  if (!areaName) {
    areaName = await proceedNewArea();
  } else {
    userName = findUserName(areaName);
  }

  if (!userName) {
    userName = await proceedNewUser(areaName)
  }

  await establish(areaName, userName);
  await run();
}

export default main;
