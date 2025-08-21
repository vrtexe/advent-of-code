function main() {
  part2();
}

export async function part1() {
  const data = await parse();
  const regex = /(?=(SAMX))|(?=(XMAS))/g;
  const [rows, columns] = [data.length, data[0].length];

  let result = 0;

  for (let row = 0; row < rows; row++) {
    result += data[row].join('').match(regex)?.length ?? 0;
  }

  for (let column = 0; column < columns; column++) {
    const columnData: string[] = [];

    for (let row = 0; row < rows; row++) {
      columnData.push(data[row][column]);
    }

    result += columnData.join('').match(regex)?.length ?? 0;
  }

  const diagonals = [
    ...extractDiagonal(rows, columns),
    ...extractDiagonalFlipped(rows, columns),
  ].map((line) => line.map(([row, column]) => data[row][column]));

  for (const diagonal of diagonals) {
    result += diagonal.join('').match(regex)?.length ?? 0;
  }

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const [rows, columns] = [data.length, data[0].length];

  let result = 0;
  for (let row = 1; row < rows - 1; row++) {
    for (let column = 1; column < columns - 1; column++) {
      if (isXmas(row, column, data)) ++result;
    }
  }

  console.log('Result:', result);
}

function isXmas(row: number, column: number, data: string[][]) {
  return (
    data[row][column] === 'A' &&
    isDiagonalXmas([row - 1, column - 1], [row + 1, column + 1], data) &&
    isDiagonalXmas([row + 1, column - 1], [row - 1, column + 1], data)
  );
}

function isDiagonalXmas(
  left: [number, number],
  right: [number, number],
  data: string[][],
) {
  const letters = new Set(['M', 'S']);

  const [leftRow, leftColumn] = left;
  const [rightRow, rightColumn] = right;

  return (
    letters.delete(data[leftRow][leftColumn]) &&
    letters.delete(data[rightRow][rightColumn])
  );
}

function extractDiagonal(rows: number, columns: number) {
  const result: [number, number][][] = [];

  for (let column = 0; column < columns; column++) {
    result.push(traverseDiagonal(1, [0, column], rows, columns));
  }

  for (let row = 1; row < rows; row++) {
    result.push(traverseDiagonal(row + 1, [row, 0], rows, columns));
  }

  return result;
}

function extractDiagonalFlipped(rows: number, columns: number) {
  const result: [number, number][][] = [];

  for (let column = 0; column < columns; column++) {
    result.push(traverseDiagonalReverse(1, [0, column], rows));
  }

  for (let row = 1; row < rows; row++) {
    result.push(traverseDiagonalReverse(row + 1, [row, columns - 1], rows));
  }

  return result;
}

function traverseDiagonal(
  index: number,
  start: [number, number],
  rows: number,
  columns: number,
) {
  const list: [number, number][] = [start];

  for (let row = index; row < rows; row++) {
    const [, lastColumn] = list.at(-1)!;
    if (lastColumn + 1 >= columns) break;
    list.push([row, lastColumn + 1]);
  }

  return list;
}

function traverseDiagonalReverse(
  index: number,
  start: [number, number],
  rows: number,
) {
  const list: [number, number][] = [start];

  for (let row = index; row < rows; row++) {
    const [, lastColumn] = list.at(-1)!;
    if (lastColumn - 1 < 0) break;
    list.push([row, lastColumn - 1]);
  }

  return list;
}

// deno-lint-ignore no-unused-vars
function printBoard(board: string[][]) {
  for (const row of board) {
    console.log(row.join(''));
  }

  console.log();
}

// deno-lint-ignore no-unused-vars
function generateBoard(rows: number, columns: number): string[][] {
  return new Array(rows).fill([]).map(() => new Array(columns).fill('.'));
}

async function parse() {
  const text = await Deno.readTextFile('src/data/task4.txt');
  return text.split('\n').map((v) => v.split(''));
}

if (import.meta.main) {
  main();
}
