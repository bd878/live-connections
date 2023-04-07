import type TitlesList from './index';

function listRecords(this: TitlesList): TextRecord[] {
  const result: TextRecord[] = [];
  for (const value of this.container.values()) {
    result.push(value);
  }
  return result;
}

export default listRecords;
