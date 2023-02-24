
function noElementCreated(moduleName: string, ...args: any[]) {
  return create(moduleName, "no element created", args);
}

function outOfRange(moduleName: string, ...args: any[]) {
  return create(moduleName, "out of range", args);
}

function wrongDataType(moduleName: string, ...args: any[]) {
  return create(moduleName, "wrong data type", args);
}

function wrongInterface(moduleName: string, ...args: any[]) {
  return create(moduleName, "wrong interface", args);
}

function failedToGet(m: string, ...args: any[]) {
  return create(m, "failed to get", args);
}

function create(moduleName: string, message: string, ...args: any[]): Error {
  return new Error(`[${moduleName}]: ${message}: ${args}`);
}

export default {
  create,
  noElementCreated,
  outOfRange,
  failedToGet,
  wrongDataType,
  wrongInterface,
};
