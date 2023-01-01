
function debounce(func, limit = 0) {
  let last = undefined;
  return (args) => {
    if (last && (Date.now() - last) < limit) {
      log.Print('[debounce]: skip');
      return;
    }

    func(args);
    last = Date.now();
  }
}

function setUrl(url) {
  if (!url) {
    throw new Error(`[setUrl]: url arg is not defined: ${url}`);
  }

  window.history.pushState(
    {}, // non-used
    "", // legacy History API
    url
  )
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

function findUserName(areaName) {
  if (!areaName) {
    return undefined;
  }

  return localStorage.getItem(areaName);
}

function bindUserToArea(area, user) {
  return localStorage.setItem(area, user);
}

function check(checkFn, trueFn, falseFn = (() => {})) {
  return (args) => checkFn() ? trueFn(args) : falseFn();
}

export {
  debounce,
  setUrl,
  takeAreaName,
  findUserName,
  bindUserToArea,
  check,
};
