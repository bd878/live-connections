import Root from '../components/Root';
import error from '../modules/error';

function setRoot() {
  const rootEl = document.getElementById("root");
  if (!rootEl) {
    throw error.failedToGet("setRoot", "root dom element");
  }

  Root.set(rootEl as HTMLDivElement);
}

export default setRoot;
