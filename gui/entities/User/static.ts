import net from 'shared/net';
import log from 'modules/log';

async function create(areaName: AreaName) {
  const user = await net.createNewUser(areaName);
  log.Print("[User create]: user:", user);
  return user;
}

export { create };
