import http from '../../net/http';
import rtti from '../../rtti';

async function create(): Promise<AreaName> {
  const response = await http.get("/area/new");
  const areaName = await response.text();
  if (!rtti.isAreaName(areaName)) {
    throw new Error(`[create Area]: response text is not area name: ${areaName}`);
  }

  return areaName;
}

export { create };
