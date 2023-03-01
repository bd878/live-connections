/* public */
function getUid(this: Elem & Identifable): Uid {
  return this.name + "-" + this.id;
}

export default getUid;
