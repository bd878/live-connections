
function takeAreaName(path: string): string {
  if (!path) {
    return "";
  }
  const parts = path.split("/");
  if (parts.length < 2) {
    return "";
  }

  return parts[1] ? "" : parts[1];
}

export default takeAreaName;
