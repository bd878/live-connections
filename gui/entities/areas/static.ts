import http from '../../net/http';
import log from '../../modules/log';
import rtti from '../../rtti';

async function create(): Promise<AreaName> {
  const response = await http.get("/area/new");
  const areaName = await response.text();
  if (!rtti.isAreaName(areaName)) {
    throw new Error(`[create Area]: response text is not area name: ${areaName}`);
  }

  log.Print("Area create", "areaName:", areaName);
  return areaName;
}

export { create };
