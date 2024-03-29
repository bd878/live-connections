import Log from '../../modules/log';

const log = new Log("net/http");

async function get(path: string): Promise<any> {
  const response = await fetch(HTTP_PROTOCOL + BACKEND_URL + path);
  if (!response.ok) {
    throw new Error("[get]: failed to create new area");
  }

  try {
    return response;
  } catch (e: any) {
    log.Fail("error occured while retrieving response body text");
    throw new Error(e);
  }
}

async function post(path: string, options: Record<string, any>): Promise<any> {
  const response = await fetch(HTTP_PROTOCOL + BACKEND_URL + path, options);
  if (!response.ok) {
    throw new Error("[http post]: failed to create new user");
  }

  try {
    return response;
  } catch (e: any) {
    log.Fail("error occured while retrieving response body text");
    throw new Error(e);
  }
}

export default { get, post };
