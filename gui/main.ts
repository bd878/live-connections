import select from './protocol/select';
import establish from './protocol/init';
import socket from './net/socket';
import Log from './modules/log';
import users from './entities/users';
import areas from './entities/areas';
import squares from './entities/squares';
import Main from './components/Main';
import Root from './components/Root';
import Square from './components/Square';
import takeAreaName from './misc/takeAreaName';
import findUserName from './misc/findUserName';
import bindUserToArea from './misc/bindUserToArea';
import attachListeners from './listeners';
import setUrl from './misc/setUrl';
import setRoot from './misc/setRoot';
import getUid from './misc/getUid';

const log = new Log('main');

/* Waits for protocol message on socket */
async function run() {
  log.Debug("run");

  let resolve: any, reject: any;
  const p = new Promise((r, j) => {
    resolve = r;
    reject = j;
  });

  try {
    while (1) {
      const message = await socket.waitMessage();

      log.Debug("on message");

      select(message);
    }

    ;(resolve && resolve(true));
  } catch (e) {
    log.Fail("failed to run");
    ;(reject && reject(e));
  }
}

/* Applies to server for new area allocation */
async function proceedNewArea(): Promise<AreaName> {
  log.Debug("proceed new area");

  const areaName = await areas.create();
  setUrl(`/${areaName}`);
  return areaName;
}

/* Applies to server for new user registration */
async function proceedNewUser(areaName: AreaName): Promise<UserName> {
  log.Debug("proceed new user");

  const userName = await users.create(areaName);
  bindUserToArea(areaName, userName);
  return userName;
}

/* Initializes internal parts: area, user, socket, protocol etc. */
async function main() {
  log.Debug("main");

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

  log.Debug("areaName, userName:", areaName, userName);

  users.setMyName(userName);
  areas.setMyName(areaName);
  socket.init(areaName, userName);
  await establish(areaName, userName);

  setRoot();
  Main.create();
  Root.append(Main);

  attachListeners()

  await run();
}

export default main;
