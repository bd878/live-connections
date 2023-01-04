function findUserName(areaName: AreaName | undefined): UserName | undefined {
  if (!areaName) {
    return undefined;
  }

  return localStorage.getItem(areaName);
}

export default findUserName;
