import log from '../modules/log';
import socket from '../net/socket';
import debounce from '../misc/debounce';
import { makeMouseMoveMessage } from '../protocol/messages';

let squareNode: null = null;

const isMousePressed = () => squareNode !== null;

function trackMouseMove() {
  log.Print("main", "track mouse moves");

  document.addEventListener(
    'mousemove',
    debounce((event: any) => {
      socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
    }),
  );
}

function trackMousePress() {
  log.Print("main", "track mouse moves");

  document.addEventListener(
    'mousedown',
    (event: any) => {
      squareNode = null;
    }
  );

  document.addEventListener(
    'mouseup',
    (event: any) => {
      ;(isMousePressed() && (squareNode = null));
    }
  );
}

function attach() {
  trackMouseMove();
  trackMousePress();
};

export default attach;
