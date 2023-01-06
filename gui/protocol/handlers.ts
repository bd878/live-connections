import socket from '../net/socket';
import log from '../modules/log';

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Print("set token:", e));
}

function onMouseMove(e: MouseMoveEvent) {
  log.Print("[onMouseMove]: e =", e);
}

function onInitMouseCoords(e: MouseMoveEvent) {
  log.Print("[onInitMouseCoords]: e =", e);
}

function onUsersOnline(users: UsersOnlineEvent) {
  log.Print("[onUsersOnline]: users =", users);
}

export {
  onAuthOk,
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
};
