import error from '../modules/error';
import socket from '../net/socket';
import log from '../modules/log';
import area from '../components/Area';
import createCursor from '../misc/createCursor';
import { isMovable } from '../rtti';

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Print("set token:", e));
}

function onMouseMove(e: MouseMoveEvent) {
  log.Print("[onMouseMove]: e =", e);

  if (area.hasElem(e.name)) {
    const cursor = area.getElem(e.name);
    if (isMovable(cursor)) {
      cursor.move(e.xPos, e.yPos);
    } else {
      throw error.wrongInterface(cursor.name, "Movable");
    }
  } else {
    const cursor = createCursor(e.name, e.xPos, e.yPos);
    area.addElem(e.name, cursor);
  }
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
