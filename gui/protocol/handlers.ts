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
import Button from '../components/Button';
import TitlesList from '../components/TitlesList';
import Square from '../components/Square';
import UserTile from '../components/UserTile';
import TextRow from '../components/TextRow';
import diff, { defaultMapper } from '../misc/diff';
import getUid from '../misc/getUid';
import getColorFromUserName from '../misc/getColorFromUserName';
import {
  trackTextInput,
  trackAddRecord,
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
  log.Info("onAuthOk", e);
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

  log.Debug("onTextInput", e);

  const square = area.getElem(getUid(Square.cname, e.name));
  if (isContainable(square)) {
    const textarea = square.getElem(getUid(TextArea.cname, e.name));
    if (isRedrawable(textarea)) {
      textarea.redraw(e.text);
    }
  }
}

function onInitSquareCoords(e: CoordsEvent) {
  log.Info("onInitSquareCoords", e);

  const sUid = getUid(Square.cname, e.name);
  if (!area.hasElem(sUid)) {
    const square = new Square(getColorFromUserName(e.name));
    square.create(e.name);
    square.redraw();
    area
      .addChild(square)
      .addElem(sUid, square);

    const textarea = square.getElem(getUid(TextArea.cname, e.name));
    const button = square.getElem(getUid(Button.cname, e.name));
    if (users.myName() === e.name) {
      squares.setMyUid(sUid);

      if (textarea instanceof TextArea) {
        // TODO: once middleware is setup, refactor
        trackTextInput(textarea);
      } else {
        log.Warn(getUid(TextArea.cname, e.name), " not a textarea instance");
      }

      if (button instanceof Button) {
        trackAddRecord(button);
      } else {
        log.Warn(getUid(Button.cname, e.name), " not a button instance");
      }
    } else {
      ;((textarea instanceof TextArea) && textarea.turnReadonly());
      ;((button instanceof Button) && button.turnReadonly());
    }

    square.move(e.xPos, e.yPos);
  }
}

function onUsersOnline(e: UsersOnlineEvent) {
  log.Info("onUsersOnline", e);

  const diffPair = diff(users.listNames(), e.users, defaultMapper);
  const leaved = diffPair[0];
  const entered = diffPair[1];

  log.Info("leaved users:", leaved);
  log.Info("entered users:", entered);

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

    // TODO: place UserTile logic into users.redraw()
    const tUid = getUid(UserTile.cname, name);
    if (!area.hasElem(tUid)) {
      const tile = new UserTile(user.color);
      tile.create(name);
      tile.redraw();
      usersList
        .addChild(tile)
        .addElem(tUid, tile);
    }

    const cUid = getUid(Cursor.cname, name);
    if (!area.hasElem(cUid)) {
      const cursor = new Cursor(user.color);
      cursor.create(name);
      cursor.redraw();
      area
        .addChild(cursor)
        .addElem(cUid, cursor);
    }
  }
}

function onListTitles(e: TitlesListEvent) {
  log.Info("onListTitles", e);

  const sUid = getUid(Square.cname, e.name);
  if (area.hasElem(sUid)) {
    const square = area.getElem(sUid);

    if (square instanceof Square) {
      const tUid = getUid(TitlesList.cname, e.name);
      if (square.hasElem(tUid)) {
        const titlesList = square.getElem(tUid);

        if (titlesList instanceof TitlesList) {
          // TODO: unify diffList logic with usersList
          const diffPair = diff<TextRecord, number>(
            titlesList.listRecords(),
            e.records,
            (v: TextRecord) => v.createdAt
          );
          const deleted = diffPair[0];
          const added = diffPair[1];

          log.Debug("deleted:", deleted);
          log.Debug("added:", added);

          for (let i = 0; i < deleted.length; i++) {
            const id = deleted[i].createdAt;

            const rUid = getUid(TextRow.cname, `${id}`);
            ;(titlesList.hasElem(rUid) && titlesList.delElem(rUid))
          }

          for (let i = 0; i < added.length; i++) {
            const textRecord = added[i];
            const id = textRecord.createdAt;

            const rUid = getUid(TextRow.cname, `${id}`);
            if (!titlesList.hasElem(rUid)) {
              const textRow = new TextRow(textRecord);
              textRow.create(`${id}`);
              textRow.redraw();
              titlesList
                .addChild(textRow)
                .addElem(rUid, textRow);
            }
          }
        } else {
          log.Warn("not TitlesList instance");
        }
      } else {
        log.Warn("square has no TitlesList instance", tUid);
      }
    } else {
      log.Warn("not Square instance");
    }
  } else {
    log.Warn("area has no square instance", sUid);
  }
}

export {
  onAuthOk,
  onMouseMove,
  onSquareMove,
  onInitSquareCoords,
  onUsersOnline,
  onTextInput,
  onListTitles,
};
