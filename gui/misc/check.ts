const fn: Fn = () => {};

function check(
  checkFn: Fn<any, boolean>,
  trueFn: Fn,
  falseFn: Fn = fn
): Fn {
  return (args: any) => checkFn() ? trueFn(args) : falseFn();
}

export default check;
