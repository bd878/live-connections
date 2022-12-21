
// constants
const BACKEND_URL = "https://localhost:8080";

// utils
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

function takeAreaName(path) {
  if (!path) {
    return "";
  }
  const parts = path.split("/");
  if (parts.length < 2) {
    return "";
  }
  return parts[1];
}

function tryFindUserName(areaName) {
  if (!areaName) {
    return undefined;
  }

  return localStorage.getItem(areaName);
}

// messages
const SIZE_PREFIX_SIZE = 2;
const MOUSE_MOVE_MESSAGE_SIZE = 9;
const MOUSE_MOVE_TYPE = 2;
const LITTLE_ENDIANNE = 1;
const ENDIANNE = LITTLE_ENDIANNE;

function genMouseMoveMessage(x, y) {
  const buffer = new ArrayBuffer(SIZE_PREFIX_SIZE + MOUSE_MOVE_MESSAGE_SIZE)
  const dv = new DataView(buffer)
  dv.setUint16(0, MOUSE_MOVE_MESSAGE_SIZE, ENDIANNE) // +2
  dv.setInt8(SIZE_PREFIX_SIZE, MOUSE_MOVE_TYPE, ENDIANNE); // +1
  dv.setFloat32(3, x, ENDIANNE); // +4
  dv.setFloat32(7, y, ENDIANNE); // +4

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

class User {
  constructor() {
    this.area = null;
    this.name = null;
  }
}

function handleMouseMove(event) {
  socket.send(genMouseMoveMessage(event.clientX, event.clientY));
}

async function runSocket() {
  socket = new Socket();
  await socket.create("localhost:8080/ws");

  if (socket.isReady()) {
    socket.onOpen(() => appendDiv(rootEl, "Socket opened"));
    socket.onMessage((event) => { console.log('receive message: ', event.data)});
    socket.onClose((event) => event.wasClean
      ? appendDiv(rootEl, `Closed cleanly: code=${event.code} reason=${event.reason}`)
      : appendDiv(rootEl, "Connection died")
    );
    socket.onError(() => appendDiv(rootEl, "Error"));

    document.addEventListener('mousemove', debounce(handleMouseMove, 300));
  } else {
    throw Error("[init]: failed to open socket");
  }
}

async function proceedNewArea() {
  const response = await fetch(BACKEND_URL + "/area/new");
  if (!response.ok) {
    throw new Error("[proceedNewArea]: failed to create new area");
  }

  try {
    const areaName = await response.text();
    console.log("areaName:", areaName);
  } catch (e) {
    console.log("error occured while retrieving response body text");
    console.error(e);
  }
}

function proceedNewUser(areaName) {
  console.log("new user with area name:", areaName);
}

function restoreSession(areaName, userName) {
  console.log('restore session for area and user:', areaName, userName);
}

function main() {
  const areaName = takeAreaName(window.location.pathname);
  if (!areaName) {
    proceedNewArea();
    return;
  }

  const userName = tryFindUserName(areaName);
  if (!userName) {
    proceedNewUser(areaName);
    return;
  }

  restoreSession(areaName, userName);
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

  main();
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", init);
} else {
  init();
}
