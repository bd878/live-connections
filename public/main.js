
function appendDiv(parentEl, textContent) {
  const el = document.createElement("div");
  el.textContent = textContent;

  parentEl.appendChild(el);
}

function launchSocket(rootEl) {
  const socket = new WebSocket("wss://127.0.0.1:8080/join");

  socket.onopen = function () {
    appendDiv(rootEl, "Socket opened");
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
