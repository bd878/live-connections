import error from '../../modules/error';

let mainEl: HTMLDivElement | null = null;

function create(): HTMLDivElement {
  mainEl = document.createElement("div");
  return mainEl;
}

function get(): HTMLDivElement {
  if (!mainEl) {
    throw error.noElementCreated("Main get");
  }

  return mainEl;
}

export default {
  create,
  get,
};
