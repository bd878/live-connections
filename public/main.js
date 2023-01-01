import { proceedNewArea, proceedNewUser, restoreSession } from './modules/scenario.js';
import { makeMouseMoveMessage } from './modules/messages.js';
import { establishProtocol } from './modules/protocol.js';
import { debounce, takeAreaName, findUserName } from './modules/utils.js';
import log from './modules/log.js';
import User from './modules/User.js';
import Socket from './modules/Socket.js';

async function run(socket, user) {
  const handlers = {
    onAuthOk: (text) => { ;(text === "ok" && user.setToken(text)); },
    onMouseMove: (message) => { log.Print("[onMouseMove]: message =", message); },
    onInitMouseCoords: (message) => { log.Print("[onInitMouseCoords]: message =", message); },
    onUsersOnline: (users) => { log.Print("[onUsersOnline]: users =", users); },
  };

  try {
    await establishProtocol(handlers, socket, user);
  } catch (e) {
    throw new Error(e);
  }

  trackMouseEvents(socket);
}

function trackMouseEvents(s /* socket */) {
  document.addEventListener(
    'mousemove',
    debounce((event) => {
      s.send(makeMouseMoveMessage(event.clientX, event.clientY));
    }),
  );
}

async function main() {
  log.mode = 'debug';

  const socket = new Socket();
  const user = new User();

  let userName;

  let areaName = takeAreaName(window.location.pathname);
  if (!areaName) {
    areaName = await proceedNewArea();
    userName = await proceedNewUser(areaName);

    user.define(areaName, userName);

    await run(socket, user);

    return;
  }

  userName = findUserName(areaName);
  if (!userName) {
    userName = await proceedNewUser(areaName)

    user.define(areaName, userName);

    await run(socket, user);

    return;
  }

  user.define(areaName, userName);

  await run(socket, user);
  await restoreSession(areaName, userName);
}

async function init() {
  const rootEl = document.getElementById("root");
  if (!rootEl) {
    throw ReferenceError("[init]: no #root");
  }

  if (!window['WebSocket']) {
    console.error("[init]: browser does not support WebSockets");
    return;
  }

  main();
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", init);
} else {
  init();
}
