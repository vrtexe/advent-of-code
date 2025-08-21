import { customParseInt } from "../common.ts";

function main() {
  part2();
}

export async function part1() {
  const levelLists = await parse();

  let result = 0;
  for (let i = 0; i < levelLists.length; i++) {
    result += isSafe(levelLists[i]) ? 1 : 0;
  }

  console.log('Result:', result);
}

export async function part2() {
  const levelLists = await parse();

  let result = 0;
  for (let i = 0; i < levelLists.length; i++) {
    result += isSafeOffset(levelLists[i]) ? 1 : 0;
  }

  console.log('Result:', result);
}

function isSafeOffset(levelList: number[]) {
  for (let j = 0; j < levelList.length; j++) {
    const newList = [...levelList.slice(0, j), ...levelList.slice(j + 1)];
    // isSafe(levelLists[i]);
    if (isSafe(newList)) return true;
  }

  return false;
}

function isSafe(levels: number[], tolerance: number = 0) {
  let result = 0;

  let badLevels = 0;
  let levelSkipped = 0;
  for (let i = 1; i < levels.length; i++) {
    const nextValue = result + (levels[i - (1 + levelSkipped)] - levels[i]);

    if (!isLevelSafe(result, nextValue)) {
      ++levelSkipped;
      ++badLevels;

      if (badLevels > tolerance) {
        return false;
      } else {
        continue;
      }
    }

    if (levelSkipped > 0) {
      levelSkipped = 0;
    }

    result = nextValue;
  }

  return true;
}

function isLevelSafe(value: number, nextValue: number) {
  if (value < 0 && nextValue >= value) return false;
  if (value > 0 && nextValue <= value) return false;
  if (nextValue === 0 || nextValue === value) return false;
  if (Math.abs(nextValue) - Math.abs(value) > 3) return false;
  return true;
}

async function parse() {
  const text = await Deno.readTextFile('src/data/task2.txt');
  return text.split('\n').map((v) => v.split(/\s+/).map(customParseInt));
}

if (import.meta.main) {
  main();
}
