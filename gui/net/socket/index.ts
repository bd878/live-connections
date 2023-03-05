import log from '../../modules/log';
import C from './const';

/* TODO: rewrite on class to open
  multiple socket connections simultaneously */

let conn: WebSocket | null = null;

let messagesBuffer: any[] = [];

function onError(err: any) {
  console.error('[onError]: error on socket', err);
}

function onClose(event: any) {
  ;(event.wasClean
    ? log.Print('onClose', `Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Print('onClose', "Connection died")
  );
}

function init(areaName: AreaName, userName: UserName) {
  log.Print("socket", "init");

  conn = new WebSocket(SOCKET_PROTOCOL + BACKEND_URL + SOCKET_PATH + "/" + areaName + "/" + userName);

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
      const onMessage = (event: any) => { // TODO: receive messages simultaneously, add buffer
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
