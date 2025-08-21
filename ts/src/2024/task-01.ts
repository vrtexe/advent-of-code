import { sortAscending, sum } from '../common.ts';

function main() {
  part2();
}

export async function part1() {
  const [leftList, rightList] = await parse();
  const [leftSortedList, rightSortedList] = [
    leftList.toSorted(sortAscending),
    rightList.toSorted(sortAscending),
  ];
  const listLen = Math.max(leftSortedList.length, rightSortedList.length);

  const resultList = [];

  for (let i = 0; i < listLen; i++) {
    const leftRoom = leftSortedList[i];
    const rightRoom = rightSortedList[i];

    resultList.push(Math.abs(leftRoom - rightRoom));
  }

  const result = resultList.reduce(sum, 0);

  console.log('Result', result);
}

export async function part2() {
  const [leftList, rightList] = await parse();
  const rightListCountMap = buildCountMap(rightList);

  const result = leftList
    .map((v) => v * (rightListCountMap.get(v) ?? 0))
    .reduce(sum);
  console.log('Result: ', result);
}

async function parse() {
  const text = await Deno.readTextFile('src/data/task1.txt');
  return text
    .split('\n')
    .map((v) => v.split(/\s+/))
    .reduce(
      (prev, [l, r]) => [
        [...prev[0], parseInt(l)],
        [...prev[1], parseInt(r)],
      ],
      <number[][]>[[], []],
    );
}

function buildCountMap(list: number[]) {
  const map = new Map<number, number>();

  for (let i = 0; i < list.length; i++) {
    map.set(list[i], (map.get(list[i]) ?? 0) + 1);
  }

  return map;
}

if (import.meta.main) {
  main();
}
