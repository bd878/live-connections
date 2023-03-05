import { log } from '../modules/log';

function debounce(func: Fn, limit = 0): Fn<any, any> {
  let last: any = undefined;
  return (args: any) => {
    if (last && (Date.now() - last) < limit) {
      log.Debug('debounce', 'skip');
      return;
    }

    func(args);
    last = Date.now();
  }
}

export default debounce;

