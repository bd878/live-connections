import error from '../modules/error';
import socket from '../net/socket';
import Log from '../modules/log';
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

const log = new Log("handlers");

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Debug("set token:", e));
}

function onMouseMove(e: CoordsEvent) {
  log.Debug("e =", e);

  const cursor = area.getCursorElem(getUid(Cursor.cname, e.name));
  cursor.move(e.xPos, e.yPos);
}

function onSquareMove(e: CoordsEvent) {
  log.Debug("e =", e);

  const square = area.getSquareElem(getUid(Square.cname, e.name));
  square.move(e.xPos, e.yPos);
}

function onTextInput(e: TextInputEvent) {
  log.Debug("e =", e);
}

function onInitSquareCoords(e: CoordsEvent) {
  log.Debug("e =", e);

  const sUid = getUid(Square.cname, e.name);
  if (!area.hasElem(sUid)) {
    log.Debug("create square =", sUid);

    ;((users.myName() === e.name) && squares.setMyUid(sUid));

    const square = new Square();
    square.setId(e.name);
    square.create();
    square.redraw();
    area.addElem(sUid, square);

    square.move(e.xPos, e.yPos);
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
      log.Debug("create tile =", tUid);
      const tile = new UserTile(user.color);
      tile.setId(name);
      tile.create();
      tile.redraw();
      usersList.addElem(tUid, tile);
    }

    const cUid = getUid(Cursor.cname, name);
    if (!area.hasElem(cUid)) {
      log.Debug("create cursor =", cUid);

      const cursor = new Cursor(user.color);
      cursor.setId(name);
      cursor.create();
      cursor.redraw();
      area.addElem(cUid, cursor);
    }
  }
}

export {
  onAuthOk,
  onMouseMove,
  onSquareMove,
  onInitSquareCoords,
  onUsersOnline,
  onTextInput,
};
