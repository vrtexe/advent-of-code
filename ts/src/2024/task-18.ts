function main() {
  part2();
}

type Direction = (typeof Direction)[keyof typeof Direction];
const Direction = {
  Up: '^',
  Right: '>',
  Down: 'v',
  Left: '<',
} as const;

const DirectionMap: Record<Direction, Position> = {
  [Direction.Up]: { y: -1, x: 0 },
  [Direction.Down]: { y: 1, x: 0 },
  [Direction.Left]: { y: 0, x: -1 },
  [Direction.Right]: { y: 0, x: 1 },
};

export async function part1() {
  const data = await parse();

  // const [rows, cols] = [7, 7];
  // const bytes = 12

  const [rows, cols] = [71, 71];
  const bytes = 1024;

  const activeObstacles = new Set(data.slice(0, bytes).map(positionToString));
  const grid = buildGrid(activeObstacles, rows, cols);

  const startPosition: Position = { x: 0, y: 0 };
  const endPosition: Position = { x: cols - 1, y: rows - 1 };

  const distances = findShortestPath<Tile>(
    {
      position: startPosition,
      weight: 1,
    },
    (v) => {
      const nextTiles: Tile[] = [];
      for (const direction of Object.values(Direction)) {
        const nextPosition = addPosition(v.position, DirectionMap[direction]);
        const nextPositionKey = positionToString(nextPosition);
        if (!grid[nextPosition.y]?.[nextPosition.x]) continue;
        if (activeObstacles.has(nextPositionKey)) continue;

        nextTiles.push({
          position: nextPosition,
          weight: 1,
        });
      }
      return nextTiles;
    },
    (v) => positionToString(v.position),
  );

  const result = distances.get(positionToString(endPosition))?.weight;

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  const size: [number, number] = [71, 71];

  const result = findBreakingPoint(size, data);

  console.log('Result:', `${result?.x},${result?.y}`);
}

function findBreakingPoint(size: [number, number], data: Position[]) {
  let left = 0;
  let right = data.length - 1;

  const [rows, cols] = size;
  const startPosition: Position = { x: 0, y: 0 };
  const endPosition: Position = { x: cols - 1, y: rows - 1 };
  const endPositionKey = positionToString(endPosition);

  while (left < right) {
    const value = Math.floor((right + left) / 2);
    const next = value + 1;

    const distances = findDistancesForBytes(value, data, startPosition, size);

    const otherDistances = findDistancesForBytes(
      next,
      data,
      startPosition,
      size,
    );

    if (
      distances.get(endPositionKey)?.weight &&
      otherDistances.get(endPositionKey)?.weight
    ) {
      left = value;
    }

    if (
      !distances.get(endPositionKey)?.weight &&
      !otherDistances.get(endPositionKey)?.weight
    ) {
      right = value;
    }

    if (
      distances.get(endPositionKey)?.weight &&
      !otherDistances.get(endPositionKey)?.weight
    ) {
      const nextActiveObstacles = getActiveObstaclePositionSet(next, data);
      const currentActiveObstacles = getActiveObstaclePositionSet(value, data);
      const [position] = nextActiveObstacles.difference(currentActiveObstacles);
      return parsePosition(position);
    }
  }
}

function parsePosition(value: string): Position {
  const [y, x] = value.slice(1, -1).split(',');
  return { y: parseInt(y), x: parseInt(x) };
}

function findDistancesForBytes(
  value: number,
  data: Position[],
  startPosition: Position,
  [rows, cols]: [number, number],
) {
  const activeObstacles = getActiveObstaclePositionSet(value, data);
  const grid = buildGrid(activeObstacles, rows, cols);
  return getShortestPathInGrid(startPosition, activeObstacles, grid);
}

function getActiveObstaclePositionSet(bytes: number, data: Position[]) {
  return new Set(data.slice(0, bytes).map(positionToString));
}

function getShortestPathInGrid(
  startPosition: Position,
  activeObstacles: Set<string>,
  grid: string[][],
) {
  const distances = findShortestPath<Tile>(
    {
      position: startPosition,
      weight: 1,
    },
    (v) => {
      const nextTiles: Tile[] = [];
      for (const direction of Object.values(Direction)) {
        const nextPosition = addPosition(v.position, DirectionMap[direction]);
        const nextPositionKey = positionToString(nextPosition);
        if (!grid[nextPosition.y]?.[nextPosition.x]) continue;
        if (activeObstacles.has(nextPositionKey)) continue;

        nextTiles.push({
          position: nextPosition,
          weight: 1,
        });
      }
      return nextTiles;
    },
    (v) => positionToString(v.position),
  );

  return distances;
}

type Tile = Vertex & {
  position: Position;
};

function buildGrid(obstacles: Set<string>, rows: number, cols: number) {
  const grid: string[][] = [];

  for (let row = 0; row < rows; row++) {
    const column: string[] = [];

    for (let col = 0; col < cols; col++) {
      column.push(obstacles.has(positionToStringOf(row, col)) ? '#' : '.');
    }

    grid.push(column);
  }

  return grid;
}

function positionToString(position: Position) {
  return positionToStringOf(position.y, position.x);
}

function positionToStringOf(row: number, col: number) {
  return `(${row},${col})`;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task18.txt');
  return text.split('\n').map((line) => {
    const [x, y] = line.split(',');
    return <Position>{ x: parseInt(x), y: parseInt(y) };
  });
}

type Position = {
  x: number;
  y: number;
};

if (import.meta.main) {
  main();
}

function addPosition(left: Position, right: Position): Position {
  return { x: left.x + right.x, y: left.y + right.y };
}

function findShortestPath<V extends Vertex>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  const distances: Map<string, Distance<V>> = new Map([
    [keyOf(vertex), { vertex, weight: 0 }],
  ]);

  const queue = createPriorityQueue<[number, V]>(
    [[0, vertex]],
    ([left], [right]) => left - right,
  );

  while (queue.size()) {
    const [, current] = queue.shift()!;
    const currentDistance = distances.get(keyOf(current))!;

    const nextVertices: [number, V][] = [];

    for (const next of getAdjacent(current)) {
      const nextDistance = distances.get(keyOf(next));
      if (
        !nextDistance ||
        nextDistance.weight > currentDistance.weight + next.weight
      ) {
        const nextWeight = currentDistance.weight + next.weight;
        distances.set(keyOf(next), {
          vertex: next,
          through: current,
          weight: nextWeight,
        });

        nextVertices.push([nextWeight, next]);
      }
    }

    queue.push(...nextVertices);
  }

  return distances;
}

function createPriorityQueue<T>(initial: T[], comparator: Comparator<T>) {
  const queue: T[] = initial;

  return {
    array: queue,
    size: () => queue.length,
    push(...items: T[]) {
      queue.push(...items);
      queue.sort(comparator);
    },
    shift() {
      return queue.shift();
    },
  };
}

type Distance<V> = {
  vertex: V;
  through?: V;
  weight: number;
};

type Comparator<T> = (left: T, right: T) => number;

type Vertex = {
  weight: number;
};
