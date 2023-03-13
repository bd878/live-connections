import error from '../modules/error';
import socket from '../net/socket';
import Log from '../modules/log';
import area from '../components/Area';
import users from '../entities/users';
import areas from '../entities/areas';
import squares from '../entities/squares';
import usersList from '../components/UsersList';
import Cursor from '../components/Cursor';
import TextArea from '../components/TextArea';
import Square from '../components/Square';
import UserTile from '../components/UserTile';
import diff from '../misc/diff';
import getUid from '../misc/getUid';
import {
  trackTextInput,
} from '../listeners';
import {
  isContainable,
  isRedrawable,
} from '../rtti';

const log = new Log("handlers");

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  log.Debug("onAuthOk", e);
}

function onMouseMove(e: CoordsEvent) {
  const cursor = area.getElem(getUid(Cursor.cname, e.name));
  cursor.move(e.xPos, e.yPos);
}

function onSquareMove(e: CoordsEvent) {
  const square = area.getElem(getUid(Square.cname, e.name));
  square.move(e.xPos, e.yPos);
}

function onTextInput(e: TextInputEvent) {
  // TODO: const textarea = registry.getElem(getUid(TextArea.cname, e.name));

  if (users.myName() !== e.name) {
    const square = area.getElem(getUid(Square.cname, e.name));
    if (isContainable(square)) {
      const textarea = square.getElem(getUid(TextArea.cname, e.name));
      if (isRedrawable(textarea)) {
        textarea.redraw(e.text);
      }
    }
  }
}

function onInitSquareCoords(e: CoordsEvent) {
  log.Debug("onInitSquareCoords", e);

  const sUid = getUid(Square.cname, e.name);
  if (!area.hasElem(sUid)) {
    const square = new Square();
    square.create(e.name);
    square.redraw();
    area.addElem(sUid, square);

    if (users.myName() === e.name) {
      squares.setMyUid(sUid);

      // TODO: once middleware is setup, refactor
      const myTextara = square.getElem(getUid(TextArea.cname, e.name));
      trackTextInput(myTextara);
    }

    square.move(e.xPos, e.yPos);
  }
}

function onUsersOnline(e: UsersOnlineEvent) {
  log.Debug("onUsersOnline", e);

  const diffPair = diff(users.listNames(), e.users);
  const leaved = diffPair[0];
  const entered = diffPair[1];

  log.Debug("leaved users:", leaved);
  log.Debug("entered users:", entered);

  users.set(areas.myName(), e.users);

  for (let i = 0; i < leaved.length; i++) {
    const name = leaved[i];

    const cUid = getUid(Cursor.cname, name);
    const sUid = getUid(Square.cname, name);
    const tUid = getUid(UserTile.cname, name);

    ;(area.hasElem(cUid) && area.delElem(cUid));
    ;(area.hasElem(sUid) && area.delElem(sUid));
    ;(usersList.hasElem(tUid) && usersList.delElem(tUid));
  }

  for (let i = 0; i < entered.length; i++) {
    const name = entered[i];
    const user = users.getByName(name);

    const tUid = getUid(UserTile.cname, name);
    if (!area.hasElem(tUid)) {
      const tile = new UserTile(user.color);
      tile.create(name);
      tile.redraw();
      usersList.addElem(tUid, tile);
    }

    const cUid = getUid(Cursor.cname, name);
    if (!area.hasElem(cUid)) {
      const cursor = new Cursor(user.color);
      cursor.create(name);
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
