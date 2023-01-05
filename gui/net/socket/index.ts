import log from 'modules/log';
import C from './const';

/* TODO: rewrite on class to open
  multiple socket connections simultaneously */

let conn: WebSocket | null = null; /* private */

function onError(err) { /* private */
  console.error('[onError]: error on socket', err);
}

function onClose(event) { /* private */
  ;(event.wasClean
    ? log.Print(`Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Print("Connection died")
  );
}

function init() {
  conn = new WebSocket("wss://" + C.BACKEND_URL + C.SOCKET_PATH);

  conn.addEventListener('error', onError);
  conn.addEventListener('close', onClose);
}

function isReady(): boolean {
  return conn.readyState === C.OPEN;
}

function send(message: any) {
  if (!conn) {
    throw new ReferenceError("[Socket]: connection is not created");
  }

  if (conn.readyState === C.CONNECTING) {
    log.Print('[Socket send]: still in connecting state');
  } else {
    conn.send(message);
  }
}

function waitOpen(): Promise {
  if (conn) {
    if (conn.readyState === C.OPEN) {
      return Promise.resolve();
    }

    return new Promise(resolve => (
      conn.addEventListener('open', resolve)
    ));
  } else {
    return Promise.reject();
  }
}

function waitMessage(): Promise {
  if (conn) {
    return new Promise(resolve => {
      const onMessage = (event) => {
        conn.removeEventListener('message', onMessage);
        resolve(event);
      };
      conn.addEventListener('message', onMessage);
    });
  } else {
    return Promise.reject();
  }
}

export {
  waitMessage,
  waitOpen,
  send,
  isReady,
  init,
};
