import get from '../../methods/get';
import has from '../../methods/has';
import create from '../../methods/create';
import getName from '../../methods/getName';
import hasElem from '../../methods/hasElem';
import addChild from '../../methods/addChild';
import addElem from '../../methods/addElem';
import getElem from '../../methods/getElem';
import delElemChild from '../../methods/delElemChild';

class Container implements
  Creatable,
  Containable
{
  static cname: string = "Container";

  root: HTMLElement | null = null;
  id: string = '';

  container: BaseContainer;

  getName = getName;

  constructor(public CSSClass: string = '') {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;

  addChild = addChild;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElem;
  delElem = delElemChild;
}

export default Container;
