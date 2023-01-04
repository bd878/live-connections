import main from "./main";

/*
 * init module
 * initializes an app.
 * It tests an environment
 * for neccessary properties
 * and calls main program
 **/

function init() {
  if (!window['WebSocket']) {
    console.error("[init]: browser does not support WebSockets");
    return;
  }

  main();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
