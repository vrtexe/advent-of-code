import { sum } from '../common.ts';

function main() {
  part2();
}

export async function part1() {
  const data = await parse();

  let result = 0;

  for (const design of data.designs) {
    result += findDesignPattern(design, data.patterns) ? 1 : 0;
  }

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  let result = 0;

  for (const design of data.designs) {
    result += findDesignPattern(design, data.patterns);
  }

  console.log('Result:', result);
}

function findDesignPattern(
  design: string,
  patterns: string[],
  cache: Map<string, number> = new Map(),
): number {
  if (!design.length) {
    return 1;
  }

  if (cache.has(design)) {
    return cache.get(design)!;
  }

  cache.set(
    design,
    patterns
      .filter((p) => design.startsWith(p))
      .map((p) => findDesignPattern(design.slice(p.length), patterns, cache))
      .reduce(sum, 0),
  );

  return cache.get(design)!;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task19.txt');
  const [patterns, designs] = text.split('\n\n');
  return <State>{
    patterns: patterns.split(', '),
    designs: designs.split('\n'),
  };
}

type State = {
  patterns: string[];
  designs: string[];
};

if (import.meta.main) {
  main();
}
