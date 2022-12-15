
function debounce(func, limit = 0) {
  let last = undefined;
  return (args) => {
    if (last && (Date.now() - last) < limit) {
      console.log('[debounce]: skip');
      return;
    }

    func(args);
    last = Date.now();
  }
}

function appendDiv(parentEl, textContent) {
  const el = document.createElement("div");
  el.textContent = textContent;

  parentEl.appendChild(el);
}

// messages
const MOUSE_MOVE_TYPE = 2;
const LITTLE_ENDIANNE = 1;

function genMouseMoveMessage(x, y) {
  const buffer = new ArrayBuffer(9)
  const dv = new DataView(buffer)
  dv.setInt8(0, MOUSE_MOVE_TYPE, LITTLE_ENDIANNE);
  dv.setFloat32(1, x, LITTLE_ENDIANNE);
  dv.setFloat32(5, y, LITTLE_ENDIANNE);

  console.log(`(${x}, ${y})`);
  return buffer;
}

// state
let rootEl = null;
let socket = null;

class Socket {
  static CONNECTING = 0;

  static OPEN = 1;

  static TIMEOUT_OPEN = 2000; // ms

  constructor() {
    this.conn = null;
  }

  async create(addr) {
    this.conn = new WebSocket("wss://" + addr);

    for (
      let i = 0;
      i < (Socket.TIMEOUT_OPEN / 5) && this.conn.readyState !== Socket.OPEN;
      i++
    ) {
      await new Promise(resolve => setTimeout(resolve, 300));
    }
  }

  isReady() {
    return this.conn.readyState === Socket.OPEN;
  }

  send(message) {
    if (!this.conn) {
      throw new ReferenceError("[Socket]: connection is not created");
    }

    if (this.conn.readyState === Socket.CONNECTING) {
      console.log('[Socket send]: still in connecting state');
    } else {
      this.conn.send(message);
    }
  }

  onOpen(callback) {
    this.conn.addEventListener('open', callback);
  }

  onMessage(callback) {
    this.conn.addEventListener('message', callback);
  }

  onClose(callback) {
    this.conn.addEventListener('close', callback);
  }

  onError(callback) {
    this.conn.addEventListener('error', callback);
  }
}

function handleMouseMove(event) {
  socket.send(genMouseMoveMessage(event.clientX, event.clientY));
}

function main() {
  socket.onOpen(() => appendDiv(rootEl, "Socket opened"));
  socket.onMessage((event) => appendDiv(rootEl, event.data));
  socket.onClose((event) => event.wasClean
    ? appendDiv(rootEl, `Closed cleanly: code=${event.code} reason=${event.reason}`)
    : appendDiv(rootEl, "Connection died")
  );
  socket.onError(() => appendDiv(rootEl, "Error"));

  document.addEventListener('mousemove', debounce(handleMouseMove, 300));
}

async function init() {
  rootEl = document.getElementById("root");
  if (!rootEl) {
    throw ReferenceError("[init]: no #root");
  }

  if (!window['WebSocket']) {
    appendDiv("[init]: browser does not support WebSockets");
    return;
  }

  socket = new Socket();
  await socket.create("localhost:8080/ws");

  if (socket.isReady()) {
    main();
  } else {
    throw Error("[init]: failed to open socket");
  }
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", init);
} else {
  init();
}
