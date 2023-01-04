import { makeAuthUserMessage } from './messages';
import messageLoop from './loop';

async function authUser(socket, user) {
  const authMessage = await makeAuthUserMessage(user.area, user.name);
  socket.send(authMessage);
}

function closeHandler(event) {
  ;(event.wasClean
    ? log.Print(`Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Print("Connection died")
  );
}

function errorHandler(event) {
  log.Print("error =", event);
}

async function establishProtocol(handlers, socket, user) {
  // establish
  await socket.waitOpen();
  await authUser(socket, user);

  socket.pushSendGuard(user.isNotAuthed)

  // run
  if (socket.isReady()) {
    log.Print("socket is running..."); // DEBUG

    socket.onMessage((event) => messageLoop(handlers, event));
    socket.onClose(closeHandler);
    socket.onError(errorHandler);
  } else {
    throw new Error("[init]: failed to open socket");
  }
}

export {
  establishProtocol,
};
