const log = {
  mode: 'silent' | 'warn' | 'debug' = 'silent',

  _isDebug() { return this.mode === 'debug'; },
  _isSilent() { return this.mode === 'silent'; },
  _isWarn() { return this.mode === 'warn'; },

  Print(message: string, ...args: any): void {
    ;(
      (this._isDebug() || this._isWarn()) &&
      (console.log(message, ...args))
    );
  }
}

export default log;
