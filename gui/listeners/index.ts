import log from '../modules/log';
import error from '../modules/error';
import { isAccessible } from '../rtti';
import socket from '../net/socket';
import squares from '../entities/squares';
import area from '../components/Area';
import debounce from '../misc/debounce';
import {
  makeMouseMoveMessage,
  makeSquareMoveMessage,
} from '../protocol/messages';

let shiftX: number = 0;
let shiftY: number = 0;

const disableDragStart = () => false;

function initDragging(clientX: number, clientY: number, node: HTMLElement) {
  shiftX = clientX - node.getBoundingClientRect().left;
  shiftY = clientY - node.getBoundingClientRect().top;
}

function computeMoveCoords(pageX: number, pageY: number): [number, number] {
  return [pageX - shiftX, pageY - shiftY];
}

function getMySquareNode(): HTMLElement {
  if (squares.inited()) {
    const uid = squares.myUid();
    if (area.hasElem(uid)) {
      const node = area.getElem(uid);
      if (isAccessible(node)) {
        return node.get();
      } else {
        log.Warn("trackMousePress getMySquareNode", uid, "not accessible");
      }
    } else {
      log.Warn("trackMousePress getMySquareNode", "area has not my square uid:", uid)
    }
  } else {
    log.Warn("trackMousePress getMySquareNode", "squares are not inited yet");
  }

  throw error.failedToGet("trackMousePress getMySquareNode", "no my square");
}

function onMouseDown(event: any) {
  const node = getMySquareNode();
  if (node.contains(event.target)) {
    log.Debug("listeners onMouseDown", "square node contains event target");

    squares.setMyPressed();
    initDragging(event.clientX, event.clientY, node);
    node.addEventListener('dragstart', disableDragStart);
  }
}

function onMouseUp(event: any) {
  if (squares.isMyPressed()) {
    log.Debug("listeners onMouseUp", "my square is pressed");

    squares.setMyNotPressed();
    const node = getMySquareNode();
    node.removeEventListener('dragstart', disableDragStart);
  }
}

const onMouseMove = debounce((event: any) => {
  socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
  if (squares.isMyPressed()) {
    const [posX, posY] = computeMoveCoords(event.pageX, event.pageY);
    socket.send(makeSquareMoveMessage(posX, posY));
  }
});

function trackMouseMove() {
  log.Debug("main", "track mouse moves");

  document.addEventListener('mousemove', onMouseMove);
}

function trackMousePress() {
  log.Debug("main", "track mouse moves");

  document.addEventListener('mousedown', onMouseDown);
  document.addEventListener('mouseup', onMouseUp);
}

function attach() {
  trackMouseMove();
  trackMousePress();
};

export default attach;
