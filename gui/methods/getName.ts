import error from '../modules/error';

/* public */
function getName(this: Elem): string {
  const name = (this.constructor as any).cname;
  if (!name) throw error.failedToGet("getName", "name");
  if (typeof name !== 'string') throw error.wrongDataType("getName", "not a string");
  return name;
}

export default getName;
