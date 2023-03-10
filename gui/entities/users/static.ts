import http from '../../net/http';
import rtti from '../../rtti';

async function create(areaName: AreaName): Promise<UserName> {
  const response = await http.post("/join", {
    method: "POST",
    headers: {
      "Content-Type": "text/plain; charset=utf-8"
    },
    body: areaName,
  });

  const userName = await response.text();
  if (!rtti.isUserName(userName)) {
    throw new Error(`[create User]: response text is not area name: ${userName}`);
  }

  return userName;
}

export { create };
