import { join } from '@std/path/join';

function main() {
  part1();
}

const BinaryOperation: Record<string, (left: number, right: number) => number> =
  {
    OR: (left: number, right: number) => left | right,
    AND: (left: number, right: number) => left & right,
    XOR: (left: number, right: number) => left ^ right,
  };

export async function part1() {
  const [, expressions] = await parse();

  console.log(parseInt('101010', 2) + parseInt('101100', 2));

  const results: [string, number][] = [];
  for (const [k, v] of expressions) {
    if (k.startsWith('z')) {
      results.push([k, v()]);
    }
  }

  const result = parseInt(
    results
      .toSorted(
        ([leftName], [rightName]) =>
          parseInt(rightName.slice(1)) - parseInt(leftName.slice(1)),
      )
      .map(([, v]) => v.toString())
      .join(''),
    2,
  );

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task24-test.txt');
  const [stateValue, expressions] = text.split('\n\n');
  const state = parseState(stateValue);
  return <[Map<string, number>, Map<string, () => number>]>[
    state,
    parseExpression(expressions, state),
  ];
}

function parseState(value: string): Map<string, number> {
  return new Map(
    value.split('\n').map((line) => {
      const [name, value] = line.split(': ');
      return [name, parseInt(value, 2)];
    }),
  );
}

function parseExpression(value: string, state: Map<string, number>) {
  const results = new Map<string, () => number>();

  for (const line of value.split('\n')) {
    const [expression, result] = line.split(' -> ');
    const [left, operation, right] = expression.split(' ');

    results.set(result, () => {
      if (!state.has(left)) {
        state.set(left, results.get(left)!());
      }

      if (!state.has(right)) {
        state.set(right, results.get(right)!());
      }

      return BinaryOperation[operation](state.get(left)!, state.get(right)!);
    });
  }

  return results;
}

if (import.meta.main) {
  main();
}
