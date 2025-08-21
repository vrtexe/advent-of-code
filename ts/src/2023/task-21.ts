import { det, inv, multiply } from 'npm:mathjs';
import { sortAscending } from '../common.ts';

type Direction = (typeof Direction)[keyof typeof Direction];
const Direction = {
  North: 'N',
  South: 'S',
  East: 'E',
  West: 'W',
} as const;

const Tile = {
  Rock: '#',
  Plot: '.',
  Start: 'S',
};

function main() {
  part2();
}

export async function part1() {
  const data = await parse();
  const startingPosition = findStart(data);

  const [result] = walkThroughMap(startingPosition, [64], data);

  console.log('Result:', result.size);
}

export async function part2() {
  const data = await parse();
  const startingPosition = findStart(data);

  const boardSize = data.length;

  const stepsTakenToFillSingleDiamond = 65;
  const steps = 26501365;
  const target = (steps - stepsTakenToFillSingleDiamond) / boardSize;

  const depths = generateDepths(3, stepsTakenToFillSingleDiamond, boardSize);
  const depthsValues = depths.map(([, depth]) => depth);
  const [, results] = walkThroughMap(startingPosition, depthsValues, data);

  const points: [number, number][] = depths.map(([x], i) => [x, results[i]]);

  const equation = findPolynomialCramer(points);

  const result = equation(target);

  console.log('Result:', result);
}

function generateDepths(count: number, singleFill: number, boardSize: number) {
  const result: [number, number][] = [];

  for (let depth = 0; depth < count; depth++) {
    result.push([depth, singleFill + boardSize * depth]);
  }

  return result;
}

function findPolynomialCramer(points: [number, number][]) {
  const A = [
    [points[0][0] ** 2, points[0][0], 1],
    [points[1][0] ** 2, points[1][0], 1],
    [points[2][0] ** 2, points[2][0], 1],
  ];

  const B = [points[0][1], points[1][1], points[2][1]];

  const [a, b, c] = solveSystem(B, A).map(Math.round);

  return (x: number) => {
    return a * x ** 2 + b * x + c;
  };
}

function solveSystem(results: number[], matrix: number[][]) {
  const nominators: number[][][] = generateCramerNominators(results, matrix);
  const denominator = det(matrix);
  const solutions: number[] = [];

  for (const nominator of nominators) {
    solutions.push(det(nominator) / denominator);
  }

  return solutions;
}

function generateCramerNominators(column: number[], matrix: number[][]) {
  const result: number[][][] = [];

  for (let i = 0; i < matrix.length; i++) {
    result.push(replaceMatrixColumn(i, column, matrix));
  }

  return result;
}

function replaceMatrixColumn(
  replaceIndex: number,
  column: number[],
  matrix: number[][],
) {
  const result: number[][] = [];
  for (let rowIndex = 0; rowIndex < matrix.length; rowIndex++) {
    const row: number[] = [];
    for (let colIndex = 0; colIndex < matrix[rowIndex].length; colIndex++) {
      row.push(
        colIndex === replaceIndex
          ? column[rowIndex]
          : matrix[rowIndex][colIndex],
      );
    }
    result.push(row);
  }
  return result;
}

// deno-lint-ignore no-unused-vars
function findPolynomial(points: [number, number][]) {
  const A = [
    [points[0][0] ** 2, points[0][0], 1],
    [points[1][0] ** 2, points[1][0], 1],
    [points[2][0] ** 2, points[2][0], 1],
  ];

  const B = [points[0][1], points[1][1], points[2][1]];

  const [a, b, c] = multiply(inv(A), B).map(Math.round);

  return (x: number) => {
    return a * x ** 2 + b * x + c;
  };
}

function toPositionString(position: Position) {
  return `${position.x},${position.y}`;
}

function fromPositionString(value: string): Position {
  const [x, y] = value.split(',');
  return { x: parseInt(x), y: parseInt(y) };
}

function walkThroughMap(
  start: Position,
  depths: number[],
  data: string[][],
): [Set<string>, number[]] {
  const positions = new Set([toPositionString(start)]);
  const results: number[] = [];
  const sortedDepths = depths.toSorted(sortAscending);

  for (let i = 0; i < Math.max(...depths); i++) {
    if (i === sortedDepths.at(0)) {
      results.push(positions.size);
      sortedDepths.shift();
    }

    const nextPositions: Position[] = [];

    for (const position of positions) {
      for (const p of getNextPositions(fromPositionString(position), data)) {
        nextPositions.push(p);
      }
    }

    positions.clear();

    for (const p of nextPositions) {
      positions.add(toPositionString(p));
    }
  }

  results.push(positions.size);

  return [positions, results];
}

function getNextPositions(position: Position, data: string[][]) {
  const result: Position[] = [];
  for (const direction of Object.values(Direction)) {
    const nextPosition = moveInDirection(position, direction);
    const tile = data
      .at(nextPosition.y % data.length)
      ?.at(nextPosition.x % data[0].length);
    if (tile === Tile.Rock) continue;
    result.push(nextPosition);
  }
  return result;
}

function findStart(data: string[][]): Position {
  for (let i = 0; i < data.length; i++) {
    for (let j = 0; j < data.length; j++) {
      if (data[i][j] === 'S') {
        return { x: j, y: i };
      }
    }
  }
  throw new Error('Could not find starting position');
}

async function parse() {
  const text = await Deno.readTextFile('assets/2023/task21.txt');
  return text.split('\n').map((s) => s.split(''));
}

if (import.meta.main) {
  main();
}

function moveInDirection(position: Position, direction: Direction): Position {
  switch (direction) {
    case Direction.North:
      return { x: position.x, y: position.y - 1 };
    case Direction.South:
      return { x: position.x, y: position.y + 1 };
    case Direction.West:
      return { x: position.x - 1, y: position.y };
    case Direction.East:
      return { x: position.x + 1, y: position.y };
  }
}

type Position = {
  x: number;
  y: number;
};
