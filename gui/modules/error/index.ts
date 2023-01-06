
function noElementCreated(moduleName: string, ...args: any[]) {
  return create(moduleName, "no element created", args);
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
  failedToGet,
};
