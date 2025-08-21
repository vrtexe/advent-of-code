import { sortAscending, sum } from '../common.ts';

function main() {
  part1();
}

export async function part1() {
  const data = await parse();
  const result = data
    .map(
      ([l, w, h]) =>
        2 * l * w + 2 * w * h + 2 * h * l + Math.min(l * w, w * h, h * l),
    )
    .reduce(sum, 0);

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  const result = data
    .map(
      ([l, w, h]) =>
        2 * [l, w, h].toSorted(sortAscending).slice(0, 2).reduce(sum, 0) +
        l * w * h,
    )
    .reduce(sum, 0);

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('src/data/task2015d2.txt');
  return text
    .split('\n')
    .map((s) => s.split('x'))
    .map(([l, w, h]) => [parseInt(l), parseInt(w), parseInt(h)]);
}

if (import.meta.main) {
  main();
}
