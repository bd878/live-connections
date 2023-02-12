import Cursor from '../components/Cursor';

function createCursor(name: string, xPos: number, yPos: number): Cursor {
  const cursor = new Cursor();
  cursor.create();

  cursor.move(xPos, yPos);

  return cursor;
}

export default createCursor;
