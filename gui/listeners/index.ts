import Log from '../modules/log';
import error from '../modules/error';
import { isAccessible } from '../rtti';
import socket from '../net/socket';
import squares from '../entities/squares';
import area from '../components/Area';
import debounce from '../misc/debounce';
import {
  makeMouseMoveMessage,
  makeSquareMoveMessage,
  makeTextInputMessage,
} from '../protocol/messages';

const log = new Log("listeners");

const cursorSize = 32;

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

function computeMouseCoords(mouseX: number, mouseY: number): [number, number] {
  return [mouseX - cursorSize, mouseY - cursorSize];
}

function getMySquareNode(): HTMLElement {
  if (squares.inited()) {
    const uid = squares.myUid();
    if (area.hasElem(uid)) {
      const node = area.getElem(uid);
      if (isAccessible(node)) {
        return node.get();
      } else {
        log.Warn(uid, "not accessible");
      }
    } else {
      log.Warn("area has not my square uid:", uid)
    }
  } else {
    log.Warn("squares are not inited yet");
  }

  throw error.failedToGet("no my square");
}

function onMouseDown(event: any) {
  const node = getMySquareNode();
  if (node.contains(event.target)) {
    log.Debug("square node contains event target");

    squares.setMyPressed();
    initDragging(event.clientX, event.clientY, node);
    node.ondragstart = disableDragStart;
  }
}

function onMouseUp(event: any) {
  if (squares.isMyPressed()) {
    log.Debug("my square is pressed");

    squares.setMyNotPressed();
    const node = getMySquareNode();
  }
}

const onMouseMove = debounce((event: any) => {
  socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
  if (squares.isMyPressed()) {
    const [posX, posY] = computeMoveCoords(event.pageX, event.pageY);
    socket.send(makeSquareMoveMessage(posX, posY));
  }
});

const onTextAreaInput = debounce(async (event: any) => {
  const message = await makeTextInputMessage(event.target.value);
  socket.send(message);
});

function trackMouseMove() {
  document.addEventListener('mousemove', onMouseMove);
}

function trackMousePress() {
  document.addEventListener('mousedown', onMouseDown);
  document.addEventListener('mouseup', onMouseUp);
}

function trackTextInput(elem: Elem & Accessible) {
  elem.get().addEventListener("input", onTextAreaInput);
}

export {
  trackTextInput,
  trackMouseMove,
  trackMousePress,
};
