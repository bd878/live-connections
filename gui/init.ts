import main from "./main";
import log from "./modules/log";

function init() {
  // log.mode = 'debug';

  log.Print("gui", "init");

  if (!window['WebSocket']) {
    console.error("[init]: browser does not support WebSockets");
    return;
  }

  try {
    main();
  } catch (e) {
    console.error("[init]: failed to run app", e);
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
