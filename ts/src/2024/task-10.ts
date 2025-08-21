import { customParseInt, sum } from '../common.ts';

type Direction = (typeof Direction)[keyof typeof Direction];
const Direction = {
  Up: '^',
  Down: 'v',
  Left: '<',
  Right: '>',
} as const;

const DirectionMap: Record<Direction, Position> = {
  [Direction.Up]: { x: 0, y: -1 },
  [Direction.Down]: { x: 0, y: 1 },
  [Direction.Right]: { x: 1, y: 0 },
  [Direction.Left]: { x: -1, y: 0 },
} as const;

type Position = {
  y: number;
  x: number;
};

function main() {
  part2();
}

export async function part1() {
  const data = await parse();
  const startingPoints = findAllStartingPoints(data).map((p) => [p]);

  const result = startingPoints
    .map((start) => findAllPaths(start, data))
    .map(countUniqueDestinations)
    .reduce(sum, 0);

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const startingPoints = findAllStartingPoints(data).map((p) => [p]);

  const result = startingPoints
    .map((start) => findAllPaths(start, data))
    .map((destination) => destination.length)
    .reduce(sum, 0);

  console.log('Result:', result);
}

function findAllPaths(start: Position[], data: number[][]) {
  const paths: Position[][] = [start];
  const destinations: Position[][] = [];

  while (paths.length) {
    const nextPaths: Position[][] = [];

    while (paths.length) {
      const path = paths.pop()!;
      const [branchingPaths, goalReached] = climb(path, data);
      destinations.push(...goalReached);
      nextPaths.push(...branchingPaths);
    }

    paths.push(...nextPaths);
  }

  return destinations;
}

function countUniqueDestinations(destinations: Position[][]) {
  const result: Set<string> = new Set();

  for (const destination of destinations) {
    result.add(positionToString(destination.at(-1)!));
  }

  return result.size;
}

function pathToString(path: Position[]) {
  return path.map(positionToString).join('-');
}

function climb(
  path: Position[],
  data: number[][],
): [Position[][], Position[][]] {
  const branchingPaths: Position[][] = [];
  const lastPosition = path.at(-1)!;
  const currentTile = data[lastPosition.y][lastPosition.x];
  const goalReached: Position[][] = [];

  for (const direction of Object.values(Direction)) {
    const nextPosition = addPosition(lastPosition, DirectionMap[direction]);
    const nextTile = data[nextPosition.y]?.[nextPosition.x];

    if (nextTile && nextTile - currentTile === 1) {
      if (nextTile === 9) {
        goalReached.push([...path, nextPosition]);
      } else {
        branchingPaths.push([...path, nextPosition]);
      }
    }
  }

  return [branchingPaths, goalReached];
}

function findAllStartingPoints(grid: number[][]) {
  const result: Position[] = [];
  for (let y = 0; y < grid.length; y++) {
    for (let x = 0; x < grid[y].length; x++) {
      if (grid[y][x] === 0) {
        result.push({ y, x });
      }
    }
  }
  return result;
}

function addPosition(left: Position, right: Position): Position {
  return { x: left.x + right.x, y: left.y + right.y };
}

function positionFromString(value: string): Position {
  const [y, x] = value.slice(1, -1).split(',');
  return { y: parseInt(y), x: parseInt(x) };
}

function positionToString(position: Position) {
  return `(${position.y},${position.x})`;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task10.txt');
  return text.split('\n').map((line) => line.split('').map(customParseInt));
}

if (import.meta.main) {
  main();
}
