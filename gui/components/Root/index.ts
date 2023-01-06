import error from '../../modules/error';

let rootEl: HTMLDivElement | null = null

function set(el: HTMLDivElement): HTMLDivElement {
  rootEl = el;
  if (!rootEl) {
    throw error.noElementCreated("Root get");
  }

  return rootEl;
}

function has(): boolean {
  return rootEl ? true : false;
}

function get(): HTMLDivElement {
  if (!rootEl) {
    throw error.noElementCreated("Root get");
  }

  return rootEl;
}

function append(el: HTMLDivElement) {
  if (!rootEl) {
    throw error.noElementCreated("Root get");
  }

  rootEl.append(el);
}

export default {
  set,
  get,
  has,
  append,
};
