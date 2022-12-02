
function appendDiv(parentEl, textContent) {
  const el = document.createElement("div");
  el.textContent = textContent;

  parentEl.appendChild(el);
}

function launchSocket(rootEl) {
  let socket = new WebSocket("wss://127.0.0.1:8080/");

  socket.onopen = function () {
    appendDiv(rootEl, "Socket opened");

    socket.send("test");
  };

  socket.onmessage = function (event) {
    appendDiv(rootEl, `Message: ${event.data}`);
  };

  socket.onclose = function (event) {
    if (event.wasClean) {
      appendDiv(`Closed cleanly: code=${event.code} reason=${event.reason}`);
    } else {
      appendDiv("Connection died");
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
