import debounce from 'misc/debounce';
import socket from 'net/socket';

function onAuthOk(text) {
  ;(text === "ok" && user.setToken(text));
}

function onMouseMove(message) {
  log.Print("[onMouseMove]: message =", message);
}

function onInitMouseCoords(message) {
  log.Print("[onInitMouseCoords]: message =", message);
}

function onUsersOnline(users) {
  log.Print("[onUsersOnline]: users =", users);
}

function trackMouseEvents() {
  document.addEventListener(
    'mousemove',
    debounce((event) => {
      socket.send(makeMouseMoveMessage(event.clientX, event.clientY));
    }),
  );
}

function trackUserInput() {
  /**/
}

export {
  onAuthOk,
  onMouseMove,
  onInitMouseCoords,
  onUsersOnline,
};
