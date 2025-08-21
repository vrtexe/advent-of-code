import { customParseInt, sum } from '../common.ts';

function main() {
  part2();
}

export async function part1() {
  const data = await parse();

  const state = applyBlink(25, data);
  const result = state.length;

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  const state = blinkUnique(75, data);
  const result = state.values().reduce(sum, 0);

  console.log('Result:', result);
}

function applyBlink(count: number, startingValues: number[]) {
  let values = [...startingValues];

  for (let i = 0; i < count; i++) {
    values = blink([...values]);
  }

  return values;
}

function blink(values: number[]) {
  const result: number[] = [];
  for (const value of values) {
    result.push(...applyRules(value));
  }
  return result;
}

function blinkUnique(count: number, startingValues: number[]) {
  const state = new Set(startingValues);
  const counts = new Map(startingValues.map((v) => [v, 1]));

  for (let i = 0; i < count; i++) {
    const currentValues = new Map(counts);
    counts.clear();

    for (const value of state) {
      const valueCount = currentValues.get(value)!;

      for (const mappedValue of applyRules(value)) {
        counts.set(mappedValue, (counts.get(mappedValue) ?? 0) + valueCount);
      }
    }

    state.clear();
    counts.keys().forEach((s) => state.add(s));
  }

  return counts;
}

function applyRules(value: number) {
  return replaceZero(value) ?? splitEvenDigits(value) ?? multiplyByYear(value);
}

function multiplyByYear(value: number): number[] {
  return [value * 2024];
}

function splitEvenDigits(value: number): number[] | undefined {
  const digits = value.toString().split('');
  if (digits.length % 2 === 1) return;

  return [
    parseInt(digits.slice(0, digits.length / 2).join('')),
    parseInt(digits.slice(digits.length / 2, digits.length).join('')),
  ];
}

function replaceZero(value: number): number[] | undefined {
  if (value !== 0) return;

  return [1];
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task11.txt');
  return text.split(' ').map(customParseInt);
}

if (import.meta.main) {
  main();
}
