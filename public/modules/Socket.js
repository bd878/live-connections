class Socket {
  static CONNECTING = 0;

  static OPEN = 1;

  static TIMEOUT_OPEN = 2000; // ms

  constructor() {
    this.conn = null;
    this.sendGuards = []
  }

  create(addr) {
    this.conn = new WebSocket("wss://" + addr);
  }

  isReady() {
    return this.conn.readyState === Socket.OPEN;
  }

  send(message) {
    if (!this.conn) {
      throw new ReferenceError("[Socket]: connection is not created");
    }

    if (this.conn.readyState === Socket.CONNECTING) {
      log.Print('[Socket send]: still in connecting state');
    } else {
      const c = fn => fn();

      if (!this.sendGuards.some(c)) {
        this.conn.send(message);
      }
    }
  }

  pushSendGuard(fn) {
    if (typeof fn === 'function') {
      this.sendGuards.push(fn);
    }
  }

  waitOpen() {
    if (this.conn) {
      if (this.conn.readyState === Socket.OPEN) {
        return Promise.resolve();
      }

      return new Promise(resolve => (
        this.conn.addEventListener('open', resolve)
      ));
    } else {
      return Promise.reject();
    }
  }

  onOpen(callback) {
    if (this.conn) {
      this.conn.addEventListener('open', callback);
    } else {
      throw new Error("[onOpen]: conn is null");
    }
  }

  waitMessage() {
    if (this.conn) {
      return new Promise(resolve => {
        const onMessage = (event) => {
          this.conn.removeEventListener('message', onMessage);
          resolve(event);
        }
        this.conn.addEventListener('message', onMessage);
      });
    } else {
      return Promise.reject();
    }
  }

  onMessage(callback) {
    if (this.conn) {
      this.conn.addEventListener('message', callback);
    } else {
      throw new Error("[onOpen]: conn is null");
    }
  }

  onClose(callback) {
    if (this.conn) {
      this.conn.addEventListener('close', callback);
    } else {
      throw new Error("[onOpen]: conn is null");
    }
  }

  onError(callback) {
    if (this.conn) {
      this.conn.addEventListener('error', callback);
    } else {
      throw new Error("[onOpen]: conn is null");
    }
  }
}

export default Socket;
