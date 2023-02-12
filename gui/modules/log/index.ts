const log = {
  mode: 'silent',

  _isDebug() { return this.mode === 'debug'; },
  _isSilent() { return this.mode === 'silent'; },
  _isWarn() { return this.mode === 'warn'; },

  Print(entity: string, message: string, ...args: any): void {
    ;(
      (this._isDebug() || this._isWarn()) &&
      (console.log(`[${entity}]: `, message, ...args))
    );
  }
}

export default log;
