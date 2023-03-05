import main from "./main";
import Log, { setMode } from "./modules/log";

const log = new Log('gui');

function init() {
  setMode('warn');

  log.Sub("init");

  if (!window['WebSocket']) {
    log.Fail("browser does not support WebSockets");
    return;
  }

  try {
    main();
  } catch (e) {
    log.Fail("failed to run app", e);
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
