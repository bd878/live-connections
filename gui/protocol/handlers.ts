import error from '../modules/error';
import socket from '../net/socket';
import log from '../modules/log';
import area from '../components/Area';
import users from '../entities/users';
import areas from '../entities/areas';
import cursors from '../entities/cursors';
import usersList from '../components/UsersList';
import Cursor from '../components/Cursor';
import UserTile from '../components/UserTile';
import diff from '../misc/diff';

/*
 * External handlers
 */

function onAuthOk(e: AuthOkEvent) {
  ;(e.text === "ok" && log.Print("onAuthOk", "set token:", e));
}

function onMouseMove(e: MouseMoveEvent) {
  log.Print("onMouseMove", "e =", e);

  cursors.set(e.name, e.xPos, e.yPos);
  area.redraw('cursor', e.name);
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
    ;(area.hasElem(current[i]) && area.delElem(current[i]));
    ;(usersList.hasElem(current[i]) && usersList.delElem(current[i]));
  }

  for (let i = 0; i < next.length; i++) {
    if (!area.hasElem(next[i])) {
      const user = users.getByName(next[i]);

      const cursor = new Cursor(user.color);
      cursor.create();
      cursor.redraw();
      area.addElem(next[i], cursor);

      const tile = new UserTile(user.color);
      tile.create();
      tile.redraw();
      usersList.addElem(next[i], tile);
    }
  }
}

export {
  onAuthOk,
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
};
