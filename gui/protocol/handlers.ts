import error from '../modules/error';
import socket from '../net/socket';
import log from '../modules/log';
import area from '../components/Area';
import users from '../entities/users';
import areas from '../entities/areas';
import coords from '../entities/coords';
import squares from '../entities/squares';
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

function onMouseMove(e: CoordsEvent) {
  log.Print("onMouseMove", "e =", e);

  coords.set(getUid(Cursor.cname, e.name), e.xPos, e.yPos);
  area.redraw('cursor', getUid(Cursor.cname, e.name));
}

function onSquareMove(e: CoordsEvent) {
  log.Print("onSquareMove", "e =", e);
}

function onInitMouseCoords(e: CoordsEvent) {
  log.Print("onInitMouseCoords", "e =", e);

  const cUid = getUid(Cursor.cname, e.name);
  if (!area.hasElem(cUid)) {
    log.Print("create cursor", cUid);

    const user = users.getByName(e.name);
    const cursor = new Cursor(user.color);
    cursor.setId(e.name);
    cursor.create();
    cursor.redraw();
    area.addElem(cUid, cursor);
  }
}

function onInitSquareCoords(e: CoordsEvent) {
  log.Print("onInitSquareCoords", "e =", e);

  const sUid = getUid(Square.cname, e.name);
  if (!area.hasElem(sUid)) {
    log.Print("create square", sUid);

    ;((users.myName() === e.name) && squares.setMyUid(sUid));

    const square = new Square();
    square.setId(e.name);
    square.create();
    square.redraw();
    area.addElem(sUid, square);
  }
}

function onUsersOnline(e: UsersOnlineEvent) {
  const diffPair = diff(users.listNames(), e.users);
  const current = diffPair[0];
  const next = diffPair[1];

  users.set(areas.myName(), e.users);

  for (let i = 0; i < current.length; i++) {
    const name = current[i];

    const cUid = getUid(Cursor.cname, name);
    const sUid = getUid(Square.cname, name);
    const tUid = getUid(UserTile.cname, name);

    ;(area.hasElem(cUid) && area.delElem(cUid));
    ;(area.hasElem(sUid) && area.delElem(sUid));
    ;(usersList.hasElem(tUid) && usersList.delElem(tUid));
  }

  for (let i = 0; i < next.length; i++) {
    const name = next[i];
    const user = users.getByName(name);

    const tUid = getUid(UserTile.cname, name);
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
  onSquareMove,
  onInitMouseCoords,
  onInitSquareCoords,
  onUsersOnline,
};
