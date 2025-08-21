function main() {
  part1();
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

const TurnMap: Record<Direction, Direction[]> = {
  [Direction.Up]: [Direction.Left, Direction.Right],
  [Direction.Down]: [Direction.Left, Direction.Right],
  [Direction.Left]: [Direction.Up, Direction.Down],
  [Direction.Right]: [Direction.Up, Direction.Down],
};

const ReverseDirection: Record<Direction, Direction> = {
  [Direction.Up]: Direction.Down,
  [Direction.Down]: Direction.Up,
  [Direction.Left]: Direction.Right,
  [Direction.Right]: Direction.Left,
};

type Tile = Vertex & {
  position: Position;
  facing: Direction;
};

export async function part1() {
  const data = await parse();

  const start = findTilePosition('S', data);
  const end = findTilePosition('E', data);

  const distances = findShortestPath<Tile>(
    {
      position: start,
      facing: Direction.Right,
      weight: 1,
    },
    (v) => {
      if (v.position.x === end.x && v.position.y === end.y) {
        return [];
      }

      return [
        {
          position: addPosition(v.position, DirectionMap[v.facing]),
          facing: v.facing,
          weight: 1,
        },
        ...TurnMap[v.facing].map((d) => ({
          position: addPosition(v.position, DirectionMap[d]),
          facing: d,
          weight: 1001,
        })),
      ].filter((v) => data[v.position.y][v.position.x] !== '#');
    },
    (v) => `(${v.position.y},${v.position.x})`,
  );

  const values = findDistancesTo(end, distances);
  const result = Math.min(...values.map((S) => S.weight));

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  const start = findTilePosition('S', data);
  const end = findTilePosition('E', data);

  const distances = findAllShortestPaths<Tile>(
    {
      position: start,
      facing: Direction.Right,
      weight: 1,
    },
    (v) => {
      if (v.position.x === end.x && v.position.y === end.y) {
        return [];
      }

      return [
        {
          position: addPosition(v.position, DirectionMap[v.facing]),
          facing: v.facing,
          weight: 1,
        },
        ...TurnMap[v.facing].map((d) => ({
          position: addPosition(v.position, DirectionMap[d]),
          facing: d,
          weight: 1001,
        })),
        {
          position: addPosition(
            v.position,
            DirectionMap[ReverseDirection[v.facing]],
          ),
          facing: ReverseDirection[v.facing],
          weight: 2001,
        },
      ].filter((v) => data[v.position.y][v.position.x] !== '#');
    },
    (v) => tileToString(v),
  );

  const minDistance = findDistancesTo(end, distances).reduce(
    minFunc(distanceCompare),
  ).vertex;

  const state: Position[] = [];
  const queue: Tile[] = [minDistance];

  const visited = new Set<string>();

  while (queue.length) {
    const current = queue.shift()!;

    if (visited.has(tileToString(current))) continue;

    visited.add(tileToString(current));
    state.push(current.position);

    const next = distances.get(tileToString(current))?.through ?? [];

    queue.push(...next);
  }

  const result = new Set(state.map(positionToString)).size;

  console.log('Result:', result);
}

function findDistancesTo<T>(end: Position, distances: Map<string, T>): T[] {
  const values: T[] = [];
  for (const [key, value] of distances) {
    if (key.includes(`(${end.y},${end.x}`)) {
      values.push(value);
    }
  }

  return values;
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

function findAllShortestPaths<V extends Vertex>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
  printBoard?: (distances: Map<string, MultiDistance<V>>) => void,
) {
  const distances: Map<string, MultiDistance<V>> = new Map([
    [keyOf(vertex), { vertex, weight: 0, through: [] }],
  ]);

  const queue: [number, V][] = [[0, vertex]];

  while (queue.length) {
    const [, current] = queue.shift()!;
    const currentDistance = distances.get(keyOf(current))!;

    const nextVertices: [number, V][] = [];

    for (const next of getAdjacent(current)) {
      const nextDistance = distances.get(keyOf(next));
      const nextWeight = currentDistance.weight + next.weight;
      const nextKey = keyOf(next);

      if (!nextDistance || nextWeight < nextDistance.weight) {
        distances.set(nextKey, {
          vertex: next,
          through: [current],
          weight: nextWeight,
        });
        nextVertices.push([nextWeight, next]);
      } else if (nextWeight === nextDistance.weight) {
        distances.get(nextKey)?.through.push(current);
        nextVertices.push([nextWeight, next]);
      }
    }

    printBoard?.(distances);
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

type MultiDistance<V> = {
  vertex: V;
  through: V[];
  weight: number;
};

type Distance<V> = {
  vertex: V;
  through?: V;
  weight: number;
};

type Comparator<T> = (left: T, right: T) => number;

type Vertex = {
  weight: number;
};

type Position = {
  x: number;
  y: number;
};

function positionToString(position: Position) {
  return `(${position.y},${position.x})`;
}

function distanceCompare(
  left: MultiDistance<Tile>,
  right: MultiDistance<Tile>,
) {
  return left.weight - right.weight;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task16.txt');
  return text.split('\n').map((line) => line.split(''));
}

if (import.meta.main) {
  main();
}

function findTilePosition(tile: string, data: string[][]): Position {
  for (let y = 0; y < data.length; y++) {
    for (let x = 0; x < data[y].length; x++) {
      if (data[y][x] === tile) {
        return { y, x };
      }
    }
  }

  throw new Error('Tile does not exist');
}

function addPosition(left: Position, right: Position): Position {
  return { x: left.x + right.x, y: left.y + right.y };
}

function minFunc<T>(compare: (left: T, right: T) => number) {
  return (result: T | undefined, current: T) =>
    !result || compare(current, result) < 0 ? current : result;
}

function tileToString(t: Tile) {
  return `(${t.position.y},${t.position.x},${t.facing})`;
}

// deno-lint-ignore no-unused-vars
function printBoard(list: Set<string>, data: string[][]) {
  console.log(
    data
      .map((line, y) =>
        line.map((tile, x) => (list.has(`(${y},${x})`) ? 'O' : tile)).join(''),
      )
      .join('\n'),
  );
}

// deno-lint-ignore no-unused-vars
function findPaths<V extends Vertex>(
  source: V,
  destination: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  let result: Set<string>[] = [];
  const queue: [V, number, Set<string>][] = [
    [source, 0, new Set([keyOf(source)])],
  ];

  let currentShortestPath: number | undefined = 115500;

  const destinationKey = keyOf(destination);
  while (queue.length) {
    const [last, weight, path] = queue.pop()!;

    if (currentShortestPath && weight > currentShortestPath) {
      continue;
    }

    if (keyOf(last) === destinationKey) {
      if (!currentShortestPath || weight < currentShortestPath) {
        currentShortestPath = weight;
        result = [path];
      }

      if (weight === currentShortestPath) {
        result.push(path);
      }

      continue;
    }

    for (const adjacent of getAdjacent(last)) {
      const adjacentKey = keyOf(adjacent);
      if (path.has(adjacentKey)) continue;
      queue.push([
        adjacent,
        weight + adjacent.weight,
        new Set(path).add(adjacentKey),
      ]);
    }
  }

  return result;
}
