const MAX_COUNT = 5;
const DURATION_MS = 1000;

function appendDiv(parentEl, textContent) {
  const el = document.createElement("div");
  el.textContent = textContent;

  parentEl.appendChild(el);
}

function startMessenging(rootEl, socket) {
  const promises = [];
  for (let i = 0; i < MAX_COUNT; i++) {
    promises.push(new Promise(resolve => {setTimeout(() => {
      socket.send(`count: ${i}`);
      resolve(`${i}`);
    }, DURATION_MS * i);}));
  }
  promises.push(new Promise(resolve => {setTimeout(() => {
    socket.send("done");
    resolve('done');
  }, MAX_COUNT * DURATION_MS);}));

  Promise.all(promises).then(() => {
    appendDiv(rootEl, "Promises resolved");
  }).catch(() => {
    appendDiv(rootEl, 'failed');
  });
}

function launchSocket(rootEl) {
  const socket = new WebSocket("wss://127.0.0.1:8080/ws");

  socket.onopen = function () {
    appendDiv(rootEl, "Socket opened");
    startMessenging(rootEl, socket);
  };

  socket.onmessage = function (event) {
    appendDiv(rootEl, `Message: ${event.data}`);
  };

  socket.onclose = function (event) {
    if (event.wasClean) {
      appendDiv(rootEl, `Closed cleanly: code=${event.code} reason=${event.reason}`);
    } else {
      appendDiv(rootEl, "Connection died");
    }
  };

  socket.onerror = function (error) {
    appendDiv(rootEl, "Error");
  };
}

function init() {
  const rootEl = document.getElementById("root");
  if (!rootEl) {
    throw ReferenceError("no #root");
  }

  launchSocket(rootEl);
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", init);
} else {
  init();
}
