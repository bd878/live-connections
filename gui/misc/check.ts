const fn: Fn = () => {};

function check(
  checkFn: Fn<any, bool>,
  trueFn: Fn,
  falseFn: Fn = fn
): any {
  return (args: any) => checkFn() ? trueFn(args) : falseFn();
}

export default check;
