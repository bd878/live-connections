import C from './const';
import log from '../../modules/log';

async function get(path: string): Promise<any> {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + path);
  if (!response.ok) {
    throw new Error("[get]: failed to create new area");
  }

  try {
    return response;
  } catch (e: any) {
    log.Print("net get", "error occured while retrieving response body text");
    throw new Error(e);
  }
}

async function post(path: string, options: Record<string, any>): Promise<any> {
  const response = await fetch(C.PROTOCOL + C.BACKEND_URL + path, options);
  if (!response.ok) {
    throw new Error("[http post]: failed to create new user");
  }

  try {
    return response;
  } catch (e: any) {
    log.Print("net post", "error occured while retrieving response body text");
    throw new Error(e);
  }
}

export default { get, post };
