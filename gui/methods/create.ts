/* public */
function create(this: Elem): HTMLElement {
  this.root = document.createElement("div");
  return this.root;
}

export default create;
