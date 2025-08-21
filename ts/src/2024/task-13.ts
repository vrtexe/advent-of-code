import { det } from 'npm:mathjs';

const INFLATION = 10000000000000;
const Cost = {
  A: 3,
  B: 1,
} as const;

function main() {
  part2();
}

export async function part1() {
  const data = await parse();

  let result = 0;

  for (const machine of data) {
    const coefficients = buildCoefficientMatrix(machine);
    const results = buildResultMatrix(machine);
    const [a, b] = solveSystem(results, coefficients);

    if (a % 1 === 0 && b % 1 === 0) {
      result += a * Cost.A + b * Cost.B;
    }
  }

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  let result = 0;

  for (const machine of data) {
    const coefficients = buildCoefficientMatrix(machine);
    const results = buildResultMatrix(machine, INFLATION);
    const [a, b] = solveSystem(results, coefficients);

    if (a % 1 === 0 && b % 1 === 0) {
      result += a * Cost.A + b * Cost.B;
    }
  }

  console.log('Result:', result);
}

function buildCoefficientMatrix(machine: Machine) {
  return [
    [machine.a.x, machine.b.x],
    [machine.a.y, machine.b.y],
  ];
}

function buildResultMatrix(machine: Machine, inflation = 0) {
  return [machine.p.x + inflation, machine.p.y + inflation];
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

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task13.txt');
  return text.split('\n\n').map(parseMachine);
}

function parseMachine(machine: string): Machine {
  const [a, b, p] = machine.split('\n');
  return {
    a: parseButton(a),
    b: parseButton(b),
    p: parsePrize(p),
  };
}

function parseButton(button: string): Increment {
  const [xi, yi] = button.split(': ')[1].split(', ');
  const [x, y] = [parseInt(xi.split('+')[1]), parseInt(yi.split('+')[1])];

  return { x, y };
}

function parsePrize(prize: string): Increment {
  const [xi, yi] = prize.split(': ')[1].split(', ');
  const [x, y] = [parseInt(xi.split('=')[1]), parseInt(yi.split('=')[1])];
  return { x, y };
}

type Increment = {
  x: number;
  y: number;
};

type Machine = {
  a: Increment;
  b: Increment;
  p: Increment;
};

if (import.meta.main) {
  main();
}
