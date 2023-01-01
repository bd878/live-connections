import log from "./log.js";
import C from "./constants.js";
import { setUrl, bindUserToArea } from "./utils.js";

async function proceedNewArea() {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + "/area/new");
  if (!response.ok) {
    throw new Error("[proceedNewArea]: failed to create new area");
  }

  try {
    const areaName = await response.text();
    if (!areaName) {
      throw new Error("[proceedNewArea]: empty area name");
    }

    log.Print('areaName:', areaName); // DEBUG
    setUrl(`/${areaName}`);

    return areaName;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    console.error(e);
  }
}

async function proceedNewUser(areaName) {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + "/join", {
    method: "POST",
    headers: {
      "Content-Type": "text/plain; charset=utf-8"
    },
    body: areaName
  });
  if (!response.ok) {
    throw new Error("[proceedNewUser]: failed to create new user");
  }

  try {
    const userName = await response.text();
    if (!userName) {
      throw new Error("[proceedNewUser]: empty user name");
    }

    log.Print("userName:", userName); // DEBUG
    bindUserToArea(areaName, userName);

    return userName;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    throw new Error(e);
  }
}

async function restoreSession(areaName, userName) {
  log.Print("areaName, userName", areaName, userName); // DEBUG
  await new Promise(resolve => resolve());
}

export {
  proceedNewArea,
  proceedNewUser,
  restoreSession,
};
