import net from 'shared/net';
import log from 'modules/log';

async function create() {
  const area = await net.createNewArea();
  log.Print("[Area create]: area:", area);
  return area;
}

export { create };
