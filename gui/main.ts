import { makeMouseMoveMessage } from './protocol/messages';
import select from './protocol/select';
import establish from './protocol/init';
import socket from './net/socket';
import log from './modules/log';
import User from './entities/User';
import Area from './entities/Area';
import Main from './components/Main';
import Root from './components/Root';
import debounce from './misc/debounce';
import takeAreaName from './misc/takeAreaName';
import findUserName from './misc/findUserName';
import bindUserToArea from './misc/bindUserToArea';
import setUrl from './misc/setUrl';
import setRoot from './misc/setRoot';

function trackMouseEvents() {
  log.Print("main", "track mouse events");

  document.addEventListener(
    'mousemove',
    debounce((event: any) => {
      socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
    }),
  );
}

/* Waits for protocol message on socket */
async function run() {
  log.Print("main", "run");

  let resolve: any, reject: any;
  const p = new Promise((r, j) => {
    resolve = r;
    reject = j;
  });

  try {
    while (1) {
      const message = await socket.waitMessage();

      log.Print("run", "on message");

      select(message);
    }

    ;(resolve && resolve(true));
  } catch (e) {
    log.Print("run", "failed to run");
    ;(reject && reject(e));
  }
}

/* Applies to server for new area allocation */
async function proceedNewArea(): Promise<AreaName> {
  log.Print("gui", "proceed new area");

  const areaName = await Area.create();
  setUrl(`/${areaName}`);
  return areaName;
}

/* Applies to server for new user registration */
async function proceedNewUser(areaName: AreaName): Promise<UserName> {
  log.Print("gui", "proceed new user");

  const userName = await User.create(areaName);
  bindUserToArea(areaName, userName);
  return userName;
}

/* Initializes internal parts: area, user, socket, protocol etc. */
async function main() {
  log.Print("gui", "main");

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

  socket.init(areaName, userName);
  await establish(areaName, userName);

  setRoot();
  Main.create();
  Root.append(Main);

  trackMouseEvents();

  await run();
}

export default main;
