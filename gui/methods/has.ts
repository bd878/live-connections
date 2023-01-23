import error from '../modules/error';

/* public */
function has(this: Elem): boolean {
  return this.root ? true : false;
}

export default has;
