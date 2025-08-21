import { sum } from '../common.ts';

function main() {
  part1();
}

type Position = {
  row: number;
  col: number;
};

const NumberKeypad = [
  ['7', '8', '9'],
  ['4', '5', '6'],
  ['1', '2', '3'],
  [' ', '0', 'A'],
];

const DirectionMap: Record<string, Position> = {
  '<': { row: 0, col: -1 },
  '>': { row: 0, col: 1 },
  '^': { row: -1, col: 0 },
  v: { row: 1, col: 0 },
};

const NumberKeypadPosition = Object.fromEntries(
  NumberKeypad.flatMap((line, row) =>
    line.map((tile, col) => <[string, Position]>[tile, { row, col }]),
  ),
);

const DirectionKeypad = [
  [' ', '^', 'A'],
  ['<', 'v', '>'],
];

const DirectionKeypadPosition = Object.fromEntries(
  DirectionKeypad.flatMap((line, row) =>
    line.map((tile, col) => <[string, Position]>[tile, { row, col }]),
  ),
);

export async function part1() {
  const data = await parse();
  const start = NumberKeypadPosition['A'];
  const results: [number, number][] = [];

  for (const values of data) {
    let currentPosition = start;
    let firstRobotPosition = DirectionKeypadPosition['A'];
    let secondRobotPosition = DirectionKeypadPosition['A'];

    const directions: string[] = [];
    const firstRobotDirections: string[] = [];
    const secondRobotDirections: string[] = [];

    for (const value of values) {
      const nextPosition = NumberKeypadPosition[value];
      const currentDirections: string[] = [];

      const colOffset = nextPosition.col - currentPosition.col;
      if (colOffset < 0) {
        currentDirections.push(...new Array(Math.abs(colOffset)).fill('<'));
      } else if (colOffset > 0) {
        currentDirections.push(...new Array(Math.abs(colOffset)).fill('>'));
      }

      const rowOffset = nextPosition.row - currentPosition.row;
      if (rowOffset < 0) {
        currentDirections.push(...new Array(Math.abs(rowOffset)).fill('^'));
      } else if (rowOffset > 0) {
        currentDirections.push(...new Array(Math.abs(rowOffset)).fill('v'));
      }

      let position = currentPosition;
      for (const direction of currentDirections) {
        position = addPosition(position, DirectionMap[direction]);
        if (NumberKeypad[position.row][position.col] === ' ') {
          currentDirections.reverse();
          break;
        }
      }

      currentPosition = nextPosition;
      currentDirections.push('A');

      directions.push(...currentDirections);

      for (const value of currentDirections) {
        const [nextPosition, firstCurrentDirections] = findNextDirection(
          value,
          firstRobotPosition,
        );

        firstRobotPosition = nextPosition;
        firstRobotDirections.push(...firstCurrentDirections);

        for (const value of firstCurrentDirections) {
          const [nextPosition, secondCurrentDirections] = findNextDirection(
            value,
            secondRobotPosition,
          );

          secondRobotPosition = nextPosition;
          secondRobotDirections.push(...secondCurrentDirections);
        }
      }
    }

    console.log(
      [
        values.join(''),
        directions.join(''),
        firstRobotDirections.join(''),
        secondRobotDirections.join(''),
      ].join('\n'),
      '\n',
    );

    results.push([secondRobotDirections.length, parseInt(values.join(''))]);
  }

  const result = results.map(([l, r]) => l * r).reduce(sum, 0);

  console.log('Result:', result);
}

function findNextDirection(
  value: string,
  start: Position,
): [Position, string[]] {
  const nextPosition = DirectionKeypadPosition[value];
  const nextDirections: string[] = [];

  const rowOffset = nextPosition.row - start.row;
  if (rowOffset < 0) {
    nextDirections.push(...new Array(Math.abs(rowOffset)).fill('^'));
  } else if (rowOffset > 0) {
    nextDirections.push(...new Array(Math.abs(rowOffset)).fill('v'));
  }

  const colOffset = nextPosition.col - start.col;
  if (colOffset < 0) {
    nextDirections.push(...new Array(Math.abs(colOffset)).fill('<'));
  } else if (colOffset > 0) {
    nextDirections.push(...new Array(Math.abs(colOffset)).fill('>'));
  }

  let position = start;
  for (const direction of nextDirections) {
    position = addPosition(position, DirectionMap[direction]);
    if (DirectionKeypad[position.row][position.col] === ' ') {
      nextDirections.reverse();
      break;
    }
  }

  nextDirections.push('A');

  return [nextPosition, nextDirections];
}

function addPosition(left: Position, right: Position): Position {
  return { row: left.row + right.row, col: left.col + right.col };
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task21-test.txt');
  return text.split('\n').map((line) => line.split(''));
}

if (import.meta.main) {
  main();
}
