import Log from '../../modules/log';
import C from './const';

const log = new Log("net/socket");

/* TODO: rewrite on class to open
  multiple socket connections simultaneously */

let conn: WebSocket | null = null;

let messagesBuffer: any[] = [];
let resolver: ((bufs: any[]) => void) | null = null;

function onError(err: any) {
  console.error('[onError]: error on socket', err);
}

function onClose(event: any) {
  ;(event.wasClean
    ? log.Debug(`Closed cleanly: code=${event.code} reason=${event.reason}`)
    : log.Debug("Connection died")
  );
}

function onMessage(event: any) {
  messagesBuffer.push(event);

  if (resolver) {
    const buffer = messagesBuffer;
    messagesBuffer = [];
    resolver(buffer);
    resolver = null;
  }
}

function waitMessages(): Promise<any[]> {
  if (conn) {
    return new Promise(resolve => {
      if (messagesBuffer.length > 0) {
        const buffer = messagesBuffer;
        messagesBuffer = [];
        resolve(buffer);
      } else {
        resolver = resolve;
      }
    });
  } else {
    return Promise.reject();
  }
}

function init(areaName: AreaName, userName: UserName) {
  log.Debug("init");

  conn = new WebSocket(SOCKET_PROTOCOL + BACKEND_URL + SOCKET_PATH + "/" + areaName + "/" + userName);

  conn.onmessage = onMessage;
  conn.onerror = onError;
  conn.onclose = onClose;
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
    log.Debug('still in connecting state');
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

export default {
  waitMessages,
  waitOpen,
  send,
  isReady,
  init,
};
