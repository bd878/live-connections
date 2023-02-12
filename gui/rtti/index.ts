export const isAreaName = (str: string): str is AreaName => typeof str === 'string';
export const isUserName = (str: string): str is UserName => typeof str === 'string';
export const isStyleable = (elem: any): elem is Styleable => elem ? Boolean(elem.CSSClass) : false;
export const isMovable = (elem: any): elem is Moveable => elem ? typeof elem.move === 'function' : false;
export const isAccessible = (elem: any): elem is Accessible => elem ? typeof elem.get === 'function' : false;

export default {
  isAreaName,
  isUserName,
  isStyleable,
  isMovable,
  isAccessible,
};
