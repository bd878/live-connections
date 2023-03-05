import log from '../modules/log';
import { isAccessible } from '../rtti';
import socket from '../net/socket';
import squares from '../entities/squares';
import area from '../components/Area';
import debounce from '../misc/debounce';
import { makeMouseMoveMessage } from '../protocol/messages';

function onMouseDown(event: any) {
  if (squares.inited()) {
    const uid = squares.myUid();
    if (area.hasElem(uid)) {
      const node = area.getElem(uid);
      if (isAccessible(node)) {
        ;(node.get().contains(event.target) && squares.setMyPressed());
      } else {
        log.Warn("trackMousePress mousedown", uid, "not accessible");
      }
    } else {
      log.Warn("trackMousePress mousedown", "area has not my square uid:", uid)
    }
  } else {
    log.Warn("trackMousePress mousedown", "squares are not inited yet");
  }
}

function onMouseUp(event: any) {
  squares.setMyNotPressed();
}

const onMouseMove = debounce((event: any) => {
  socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
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
