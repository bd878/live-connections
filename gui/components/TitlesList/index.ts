import get from '../../methods/get';
import has from '../../methods/has';
import setId from '../../methods/setId';
import getUid from '../../methods/getUid';
import create from '../../methods/create';
import getName from '../../methods/getName';
import hasElem from '../../methods/hasElem';
import getElem from '../../methods/getElem';
import delElemChild from '../../methods/delElemChild';
import addElemChild from '../../methods/addElemChild';
import listRecords from './listRecords';

class TitlesList implements
  Elem,
  Creatable,
  Accessible,
  Styleable,
  Containable,
  Identifable
{
  static cname: string = "TitlesList";

  root: HTMLElement | null = null;
  id: string = '';
  CSSClass: string = "titles-list";

  container: BaseContainer;

  getName = getName;

  constructor() {
    this.container = new Map();
  }

  get = get;
  has = has;
  create = create;

  listRecords = listRecords;

  hasElem = hasElem;
  getElem = getElem;
  addElem = addElemChild;
  delElem = delElemChild;

  setId = setId;
  getUid = getUid;
}

export default TitlesList;
