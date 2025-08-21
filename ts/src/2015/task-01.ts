import { sum } from '../common.ts';

const FloorInstruction: Record<string, number> = {
  '(': 1,
  ')': -1,
};

function main() {
  part2();
}

export async function part1() {
  const data = await parse();
  const result = data.map((s) => FloorInstruction[s]).reduce(sum, 0);

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  let result = 0;
  let basementFirstReached = 0;
  for (
    basementFirstReached = 0;
    basementFirstReached < data.length;
    basementFirstReached++
  ) {
    result += FloorInstruction[data[basementFirstReached]];
    if (result === -1) break;
  }

  console.log('Result:', basementFirstReached + 1);
}

async function parse() {
  const text = await Deno.readTextFile('src/data/task2015.txt');
  return text.split('');
}

if (import.meta.main) {
  main();
}
