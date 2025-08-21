import { customParseInt } from '../common.ts';

type Equation = {
  result: number;
  args: number[];
};

type Operation = (typeof Operation)[keyof typeof Operation];
const Operation = {
  Add: '+',
  Multiply: '*',
  Concatenation: '||',
};

const Part1Operations: Operation[] = [Operation.Add, Operation.Multiply];
const Part2Operations: Operation[] = Object.values(Operation);

function main() {
  part2();
}

export async function part1() {
  const data = await parse();
  const result = sumPossibleEquationResults(data, Part1Operations);

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const result = sumPossibleEquationResults(data, Part2Operations);

  console.log('Result:', result);
}

function sumPossibleEquationResults(
  equations: Equation[],
  operations: Operation[],
) {
  let result = 0;

  for (const equation of equations) {
    if (isPossible(equation, operations)) {
      result += equation.result;
    }
  }

  return result;
}

function isPossible(equation: Equation, operations: Operation[]) {
  const results = calculateAllResults(equation, operations);
  return results.has(equation.result);
}

function calculateAllResults(equation: Equation, operations: Operation[]) {
  const [first, ...args] = equation.args;
  let results = [first];

  for (const arg of args) {
    results = expandResults(arg, results, equation.result, operations);
  }

  return new Set(results);
}

function expandResults(
  next: number,
  results: number[],
  expected: number,
  operations: Operation[],
) {
  const newResults = [];

  for (const result of results) {
    for (const nextResult of expandByOperation(result, next, operations)) {
      if (nextResult > expected) {
        continue;
      }
      newResults.push(nextResult);
    }
  }

  return newResults;
}

function expandByOperation(
  left: number,
  right: number,
  operations: Operation[],
): number[] {
  const results: number[] = [];

  for (const operation of operations) {
    results.push(performOperation(left, right, operation));
  }

  return results;
}

function performOperation(left: number, right: number, operation: Operation) {
  switch (operation) {
    case Operation.Add:
      return left + right;
    case Operation.Multiply:
      return left * right;
    case Operation.Concatenation:
      return parseInt(`${left}${right}`);
  }

  return 0;
}

async function parse(): Promise<Equation[]> {
  const text = await Deno.readTextFile('src/data/task7.txt');
  return text.split('\n').map((equation) => {
    const [result, args] = equation.split(/:\s+/);
    return {
      result: parseInt(result),
      args: args.split(/\s+/).map(customParseInt),
    };
  });
}

if (import.meta.main) {
  main();
}
