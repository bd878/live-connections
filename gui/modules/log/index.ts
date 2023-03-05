type LogMode = 'silent' | 'debug' | 'warn';

let mode: LogMode = 'silent';

const isDebug = (): boolean => { return mode === 'debug'; };
const isSilent = (): boolean => { return mode === 'silent'; };
const isWarn = (): boolean => { return mode === 'warn'; };

const setMode = (m: LogMode) => { mode = m; };

class Log {
  sub: string = '';

  constructor(public prefix: string = '') {}

  Sub(sub: string) {
    this.sub = sub;
  }

  Debug(entity: string, ...args: any): void {
    if (!isSilent() && (isDebug() || isWarn())) {
      this.print(entity, ...args);
    }
  }

  Fail(entity: string, ...args: any): void {
    if (!isSilent()) {
      this.print(entity, ...args);
    }
  }

  Warn(entity: string, ...args: any): void {
    if (!isSilent() && isWarn()) {
      this.print(entity, ...args);
    }
  }

  print(entity: string, ...args: any): void {
    let start: string = '';

    if (this.prefix) {
      start = `[${this.prefix}`;
    } else if (entity) {
      start = `[${entity}`;
    }

    if (this.sub) {
      start += ` ${this.sub}]`;
    } else if (start) {
      start += ']';
    }

    console.log(start, entity, ...args);
  }
}

const log = new Log();

export default Log;
export { log, setMode };
