const log = {
  mode: 'silent',

  _isDebug() { return this.mode === 'debug'; },
  _isSilent() { return this.mode === 'silent'; },
  _isWarn() { return this.mode === 'warn'; },

  Debug(entity: string, message: string, ...args: any): void {
    log.Print(entity, message, ...args);
  },

  Fail(entity: string, message: string, ...args: any): void {
    log.Print(entity, message, ...args);
  },

  Warn(entity: string, message: string, ...args: any): void {
    log.Print(entity, message, ...args);
  },

  Print(entity: string, message: string, ...args: any): void {
    ;(
      (this._isDebug() || this._isWarn()) &&
      (console.log(`[${entity}]:`, message, ...args))
    );
  }
}

export default log;
