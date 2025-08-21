import { customParseInt, sum } from '../common.ts';

function main() {
  part2();

  console.log()
}

export async function part1() {
  const data = await parse();

  let endIndex = data.length - 1;
  let endIndexCount = 0;
  let endCount = 0;

  let result = 0;
  let totalIndex = 0;

  for (let i = 0; i < data.length - endCount; i++) {
    if (i >= endIndex) {
      for (let j = endIndexCount; j < data[i]; j++) {
        result += totalIndex++ * (endIndex / 2);
      }
      break;
    }

    if (i % 2 === 0) {
      for (let j = 0; j < data[i]; j++) {
        result += totalIndex++ * (i / 2);
      }
    } else {
      for (let j = 0; j < data[i]; j++) {
        if (endIndexCount >= data[endIndex]) {
          endIndex -= 2;
          endIndexCount = 0;
          ++endCount;
        }
        result += totalIndex++ * (endIndex / 2);
        ++endIndexCount;
      }
    }
  }

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const state: [number[], number][] = [];

  for (let count = 0; count < data.length; count++) {
    if (count % 2 === 0) {
      state.push([new Array(data[count]).fill(count / 2), 0]);
    } else {
      state.push([[], data[count]]);
    }
  }

  let fullSpaces = 1;
  let orderedFullSpaces = 0;

  for (let even = data.length - 1; even >= 0; even -= 2) {
    const [evenData] = state[even];

    let inOrder = false;
    for (let space = fullSpaces; space < even; space++) {
      const [spaceData, spaceCount] = state[space];

      if (spaceCount === 0 && inOrder) {
        orderedFullSpaces++;
      } else {
        inOrder = false;
      }

      if (evenData.length <= spaceCount) {
        spaceData.push(...evenData);
        state[space] = [spaceData, spaceCount - evenData.length];
        state[even] = [[], evenData.length];
        break;
      }
    }

    fullSpaces = orderedFullSpaces;
    orderedFullSpaces = 0;
  }

  const result = state
    .flatMap(([d, c]) => <number[]>[...d, ...new Array(c).fill(0)])
    .map((v, i) => i * v)
    .reduce(sum, 0);

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task9.txt');
  return text.split('').map(customParseInt);
}

if (import.meta.main) {
  main();
}
