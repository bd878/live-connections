import socket from '../net/socket';
import log from '../modules/log';

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Print("set token:", e));
}

function onMouseMove(e: MouseMoveEvent) {
  log.Print("[onMouseMove]: e =", e);
}

function onInitMouseCoords(e: MouseMoveEvent) {
  log.Print("[onInitMouseCoords]: e =", e);
}

function onUsersOnline(e: UsersOnlineEvent) {
  log.Print("[onUsersOnline]: users =", e.users);
}

export {
  onAuthOk,
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
};
