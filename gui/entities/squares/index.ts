import error from '../../modules/error';

let _inited: boolean = false;
let pressed: boolean = false;
let square: Uid | null = null;

function myUid(): Uid {
  if (!square) throw error.noElementCreated("squares", "my");
  return square;
}

function setMyUid(uid: Uid) {
  square = uid;
  _inited = true;
}

const inited = (): boolean => _inited;

const isMyPressed = (): boolean => pressed;
const setMyPressed = (): void => {pressed = true;};
const setMyNotPressed = (): void => {pressed = false;};

export default {
  inited,
  myUid,
  isMyPressed,
  setMyPressed,
  setMyNotPressed,
  setMyUid,
};
