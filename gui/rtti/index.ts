export const isAreaName = (str: string): str is AreaName => typeof str === 'string';
export const isUserName = (str: string): str is UserName => typeof str === 'string';
export const isStyleable = (elem: any): elem is Styleable => !!elem.CSSClass;

export default {
  isAreaName,
  isUserName,
  isStyleable,
};
