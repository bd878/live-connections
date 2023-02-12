import log from '../../modules/log';
import C from './const';

/* TODO: rewrite on class to open
  multiple socket connections simultaneously */

let conn: WebSocket | null = null; /* private */

function onError(err: any) { /* private */
  console.error('[onError]: error on socket', err);
}

function onClose(event: any) { /* private */
  ;(event.wasClean
    ? log.Print('onClose', `Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Print('onClose', "Connection died")
  );
}

function init() {
  log.Print("socket", "init");

  conn = new WebSocket(C.PROTOCOL + C.BACKEND_URL + C.SOCKET_PATH);

  conn.addEventListener('error', onError);
  conn.addEventListener('close', onClose);
}

function isReady(): boolean {
  if (!conn) {
    throw new ReferenceError("[Socket]: connection is not created");
  }

  return conn.readyState === C.OPEN;
}

function send(message: any): void {
  if (!conn) {
    throw new ReferenceError("[Socket]: connection is not created");
  }

  if (conn.readyState === C.CONNECTING) {
    log.Print('Socket send', 'still in connecting state');
  } else {
    conn.send(message);
  }
}

function waitOpen(): Promise<any> {
  if (conn) {
    if (conn.readyState === C.OPEN) {
      return Promise.resolve();
    }

    return new Promise(resolve => {
      ;(conn && conn.addEventListener('open', resolve));
    });
  } else {
    return Promise.reject();
  }
}

function waitMessage(): Promise<void> {
  if (conn) {
    return new Promise(resolve => {
      const onMessage = (event: any) => {
        ;(conn && conn.removeEventListener('message', onMessage));
        resolve(event);
      };
      ;(conn && conn.addEventListener('message', onMessage));
    });
  } else {
    return Promise.reject();
  }
}

export default {
  waitMessage,
  waitOpen,
  send,
  isReady,
  init,
};
