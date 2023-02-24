import error from '../modules/error';

const COLOR_LENGTH = 6;

function getColorFromUserName(name: UserName): Color {
  const color = name.slice(0, COLOR_LENGTH); // hex
  if (color.length < COLOR_LENGTH) {
    throw error.wrongDataType("getColorFromUserName", "not enough user name length")
  }
  return `#${color}`;
}

export default getColorFromUserName;