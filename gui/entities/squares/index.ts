import error from '../../modules/error';

let square: Uid | null = null;

function myUid(): Uid {
  if (!square) throw error.noElementCreated("squares", "my");
  return square;
}

function setMyUid(uid: Uid) {
  square = uid;
}

export default {
  myUid,
  setMyUid,
};
