function bindUserToArea(areaName: AreaName, userName: UserName) {
  localStorage.setItem(areaName, userName);
}

export default bindUserToArea;
