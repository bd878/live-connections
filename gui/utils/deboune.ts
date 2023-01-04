import log from 'modules/log';

function debounce(func: Fn, limit = 0) {
  let last = undefined;
  return (args) => {
    if (last && (Date.now() - last) < limit) {
      log.Print('[debounce]: skip');
      return;
    }

    func(args);
    last = Date.now();
  }
}

export default debounce;

