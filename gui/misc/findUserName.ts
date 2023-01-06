function findUserName(areaName: AreaName | undefined) {
  if (!areaName) {
    return undefined;
  }

  return localStorage.getItem(areaName);
}

export default findUserName;
