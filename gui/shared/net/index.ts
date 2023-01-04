import C from 'constants';

async function createNewArea(): AreaName {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + "/area/new");
  if (!response.ok) {
    throw new Error("[proceedNewArea]: failed to create new area");
  }

  try {
    const areaName = await response.text();
    if (!areaName) {
      throw new Error("[proceedNewArea]: empty area name");
    }

    return areaName as AreaName;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    throw new Error(e);
  }
}

async function createNewUser(areaName: AreaName): UserName {
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

    return userName as UserName;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    throw new Error(e);
  }
}

export { createNewUser, createNewArea };
