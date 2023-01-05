import C from './const';

async function get(path: string) {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + path);
  if (!response.ok) {
    throw new Error("[get]: failed to create new area");
  }

  try {
    await response.blob(); // try - fail

    return response;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    throw new Error(e);
  }
}

async function post(path: string, options: Record<string, any>) {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + path, options);
  if (!response.ok) {
    throw new Error("[http post]: failed to create new user");
  }

  try {
    await response.blob(); // try - fail

    return response;
  } catch (e) {
    log.Print("error occured while retrieving response body text");
    throw new Error(e);
  }
}

export { get, post };
