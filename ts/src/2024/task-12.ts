import { customParseInt } from '../common.ts';

function main() {
  part2();
}

type WallType = (typeof WallType)[keyof typeof WallType];
const WallType = {
  Vertical: '|',
  Horizontal: '-',
} as const;

const WallMap: Record<string, WallType> = Object.fromEntries(
  Object.entries(WallType).map(([_, v]) => [v, v]),
);

const Walls: string[] = [WallType.Vertical, WallType.Horizontal];

const Directions: [number, number, WallType][] = [
  [0, -1, WallType.Vertical],
  [0, 1, WallType.Vertical],
  [-1, 0, WallType.Horizontal],
  [1, 0, WallType.Horizontal],
];

const VerticalDirections: [number, number][] = [
  [-1, 0],
  [1, 0],
];

const HorizontalDirections: [number, number][] = [
  [0, -1],
  [0, 1],
];

const WallDirection: Record<WallType, [number, number][]> = {
  [WallType.Horizontal]: HorizontalDirections,
  [WallType.Vertical]: VerticalDirections,
};

export async function part1() {
  const data = await parse();
  const expandedGrid: string[][] = fillExpandedGrid(expandGrid(data));
  const regions = findSections(expandedGrid);
  const result = calculatePrice(regions);

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const expandedGrid: string[][] = fillExpandedGrid(expandGrid(data));
  const regions = findSections(expandedGrid);
  const result = calculateDiscountedPrice(regions);

  // 936718
  console.log('Result:', result);
}

type RegionData = [number, number, number, string];

function calculatePrice(regions: RegionData[]) {
  let result = 0;

  for (const [area, perimeter] of regions) {
    result += area * perimeter;
  }

  return result;
}

function calculateDiscountedPrice(regions: RegionData[]) {
  let result = 0;

  for (const [area, , sides] of regions) {
    result += area * sides;
  }

  return result;
}

function findSections(grid: string[][]) {
  const passed: Set<string> = new Set();
  const section: RegionData[] = [];

  const start: [number, number] = [1, 1];
  const nextPositions: [number, number][] = [start];
  const nextUniquePositions: Set<string> = new Set([keyOfPosition(...start)]);

  while (nextPositions.length) {
    const [startRow, startCol] = nextPositions.pop()!;
    const startKey = keyOfPosition(startRow, startCol);
    nextUniquePositions.delete(startKey);

    if (passed.has(keyOfPosition(startRow, startCol))) continue;

    const region: [number, number][] = [[startRow, startCol]];

    let regionPerimeter = 0;
    let regionArea = 0;
    const walls: Set<string> = new Set();

    while (region.length) {
      const [row, col] = region.pop()!;
      const visited: Set<string> = new Set();

      if (!passed.has(keyOfPosition(row, col))) {
        passed.add(keyOfPosition(row, col));
        ++regionArea;
      } else continue;

      // printGrid(passed, grid);

      for (const [r, c] of Directions) {
        const [nextRow, nextCol] = [row + r, col + c];
        const wallKey = keyOfPosition(nextRow, nextCol);
        if (!visited.has(wallKey)) {
          visited.add(wallKey);
        } else continue;

        const regionPosition: [number, number] = [row + r * 2, col + c * 2];
        const [regionRow, regionCol] = regionPosition;
        const regionPositionKey = keyOfPosition(regionRow, regionCol);

        if (Walls.includes(grid[nextRow][nextCol])) {
          if (
            grid[regionRow]?.[regionCol] &&
            !nextUniquePositions.has(regionPositionKey) &&
            !passed.has(regionPositionKey)
          ) {
            nextPositions.push([regionRow, regionCol]);
            nextUniquePositions.add(regionPositionKey);
          }

          walls.add(wallKey);
          ++regionPerimeter;
          continue;
        }

        if (passed.has(regionPositionKey)) continue;
        region.push(regionPosition);
      }
    }

    section.push([
      regionArea,
      regionPerimeter,
      countSides(walls, grid),
      grid[startRow][startCol],
    ]);
  }

  return section;
}

function keyOfPosition(row: number, col: number) {
  return `${row},${col}`;
}

function fillExpandedGrid(expandedGrid: string[][]) {
  for (let row = 1; row < expandedGrid.length - 1; row += 2) {
    for (let col = 1; col < expandedGrid[row].length - 1; col += 2) {
      const current = expandedGrid[row][col];
      for (const [r, c, i] of Directions) {
        const next = expandedGrid[row + r * 2]?.[col + c * 2];
        if (!next || current !== next) {
          expandedGrid[row + r][col + c] = i;
        }
      }
    }
  }

  return expandedGrid;
}

function expandGrid(grid: string[][]) {
  const emptyExpansionRow = generateExpansionRow(grid[0].length);
  const expandedInnerGrid = expandInnerGrid(grid);
  const expandedGrid = `${emptyExpansionRow}\n${expandedInnerGrid}\n${emptyExpansionRow}`;
  return parseText(expandedGrid);
}

function expandInnerGrid(grid: string[][]) {
  const emptyExpansionRow = generateExpansionRow(grid[0].length);
  return grid
    .map((line) => ` ${line.join(' ')} `)
    .join(`\n${emptyExpansionRow}\n`);
}

function generateExpansionRow(size: number) {
  return ` ${new Array(size).fill(' ').join(' ')} `;
}

function parsePosition(value: string): [number, number] {
  return value.split(',').map(customParseInt) as [number, number];
}

function countSides(walls: Set<string>, grid: string[][]) {
  let result = 0;
  let nextWalls = new Set(walls);

  while (nextWalls.size) {
    const [side, ...otherSides] = nextWalls;
    nextWalls = new Set(otherSides);

    const sameSideWalls = findWallsOnSide(parsePosition(side), nextWalls, grid);
    for (const element of sameSideWalls) {
      nextWalls.delete(element);
    }

    ++result;
  }

  return result;
}

function findWallsOnSide(
  [row, col]: [number, number],
  walls: Set<string>,
  grid: string[][],
) {
  const direction = WallMap[grid[row][col]];
  const oppositeDirection = getOppositeDirection(direction);
  const result: string[] = [];

  for (const [r, c] of WallDirection[direction]) {
    let [nextRow, nextCol] = [row + r * 2, col + c * 2];
    while (
      walls.has(keyOfPosition(nextRow, nextCol)) &&
      grid[nextRow]?.[nextCol] === direction &&
      isWallCross([nextRow - r, nextCol - c], oppositeDirection, grid)
    ) {
      result.push(keyOfPosition(nextRow, nextCol));
      [nextRow, nextCol] = [nextRow + r * 2, nextCol + c * 2];
    }
  }

  return result;
}

function isWallCross(
  [row, col]: [number, number],
  direction: WallType,
  grid: string[][],
) {
  const [[prevRow, prevCol], [nextRow, nextCol]] = WallDirection[direction];
  return (
    grid[row + prevRow]?.[col + prevCol] !== direction ||
    grid[row + nextRow]?.[col + nextCol] !== direction
  );
}

function getOppositeDirection(direction: WallType) {
  switch (direction) {
    case WallType.Horizontal:
      return WallType.Vertical;
    case WallType.Vertical:
      return WallType.Horizontal;
  }
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task12.txt');
  return parseText(text);
}

function parseText(text: string) {
  return text.split('\n').map((line) => line.split(''));
}

// deno-lint-ignore no-unused-vars
function printGrid(marked: Set<string>, grid: string[][]) {
  console.log(
    grid
      .map((s, r) =>
        s.map((d, c) => (marked.has(keyOfPosition(r, c)) ? '+' : d)),
      )
      .map((s) => s.join(''))
      .join('\n'),
    '\n',
  );
}

if (import.meta.main) {
  main();
}
