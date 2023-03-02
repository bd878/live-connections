import error from '../modules/error';
import socket from '../net/socket';
import log from '../modules/log';
import area from '../components/Area';
import users from '../entities/users';
import areas from '../entities/areas';
import cursors from '../entities/cursors';
import usersList from '../components/UsersList';
import Cursor from '../components/Cursor';
import Square from '../components/Square';
import UserTile from '../components/UserTile';
import diff from '../misc/diff';
import getUid from '../misc/getUid';

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Print("onAuthOk", "set token:", e));
}

function onMouseMove(e: MouseMoveEvent) {
  log.Print("onMouseMove", "e =", e);

  cursors.set(getUid(Cursor.name, e.name), e.xPos, e.yPos);
  area.redraw('cursor', getUid(Cursor.name, e.name));
}

function onInitMouseCoords(e: MouseMoveEvent) {
  log.Print("onInitMouseCoords", "e =", e);
}

function onUsersOnline(e: UsersOnlineEvent) {
  const diffPair = diff(users.listNames(), e.users);
  const current = diffPair[0];
  const next = diffPair[1];

  users.set(areas.my().name, e.users);

  for (let i = 0; i < current.length; i++) {
    const name = current[i];

    const cUid = getUid(Cursor.name, name);
    const sUid = getUid(Square.name, name);
    const tUid = getUid(UserTile.name, name);

    ;(area.hasElem(cUid) && area.delElem(cUid));
    ;(area.hasElem(sUid) && area.delElem(sUid));
    ;(usersList.hasElem(tUid) && usersList.delElem(tUid));
  }

  for (let i = 0; i < next.length; i++) {
    const name = next[i];
    const user = users.getByName(name);

    const cUid = getUid(Cursor.name, name);
    if (!area.hasElem(cUid)) {
      log.Print("create cursor", cUid);
      const cursor = new Cursor(user.color);
      cursor.setId(name);
      cursor.create();
      cursor.redraw();
      area.addElem(cUid, cursor);
    }

    const sUid = getUid(Square.name, name);
    if (!area.hasElem(sUid)) {
      log.Print("create square", sUid);
      const square = new Square();
      square.setId(name);
      square.create();
      square.redraw();
      area.addElem(sUid, square);
    }

    const tUid = getUid(UserTile.name, name);
    if (!area.hasElem(tUid)) {
      log.Print("create tile", tUid);
      const tile = new UserTile(user.color);
      tile.setId(name);
      tile.create();
      tile.redraw();
      usersList.addElem(tUid, tile);
    }
  }
}

export {
  onAuthOk,
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
};
