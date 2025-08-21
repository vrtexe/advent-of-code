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

function main() {
  part1();
}

export async function part1() {
  const data = await parse();

  const startPosition = findTilePosition('S', data);
  const endPosition = findTilePosition('E', data);

  const distances = findShortestPathFor(startPosition, data);
  const distance = distances.get(positionToString(endPosition))!;

  // for (const node of distance.path) {
  const cheats = new Map<number, number>();
  for (let nodeIndex = 0; nodeIndex < distance.path.length - 2; nodeIndex++) {
    const node = distance.path[nodeIndex];
    let shortcut = 0;
    for (let i = nodeIndex + 2; i < distance.path.length; i++) {
      const next = distance.path[i];
      if (positionDistance(node.position, next.position) === 3) {
        shortcut = i;
      }
    }
    const savedTime = shortcut - nodeIndex - 3;

    // if (shortcut <= 2) continue;

    cheats.set(savedTime, (cheats.get(savedTime) ?? 0) + 1);
    console.log(
      `Can cut: ${shortcut - nodeIndex - 1}ps (${nodeIndex}, ${shortcut})`,
    );

    // break;
  }

  const pathMap = new Map(
    distance.path.map((t, i) => [
      positionToString(t.position),
      i.toString(),
    ]) as [string, string][],
  );
  printBoard(pathMap, data);

  console.log(cheats);
  // }

  // const result = data.length;
  console.log('Result:', startPosition);
}

function positionDistance(left: Position, right: Position) {
  return Math.abs(right.x - left.x) + Math.abs(right.y - left.y);
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task20-test.txt');
  return text.split('\n').map((line) => line.split(''));
}

if (import.meta.main) {
  main();
}

function findShortestPathFor(startPosition: Position, grid: string[][]) {
  return findShortestPath<Tile>(
    {
      position: startPosition,
      weight: 1,
    },
    (v) => {
      const nextTiles: Tile[] = [];
      for (const direction of Object.values(Direction)) {
        const nextPosition = addPosition(v.position, DirectionMap[direction]);
        if (
          !grid[nextPosition.y]?.[nextPosition.x] ||
          grid[nextPosition.y]?.[nextPosition.x] === '#'
        )
          continue;

        nextTiles.push({
          position: nextPosition,
          weight: 1,
        });
      }
      return nextTiles;
    },
    (v) => positionToString(v.position),
  );
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

function findShortestPath<V extends Vertex>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  const distances: Map<string, Distance<V>> = new Map([
    [keyOf(vertex), { vertex, weight: 0, path: [vertex] }],
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
          path: [...currentDistance.path, next],
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
  path: V[];
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

type Tile = Vertex & {
  position: Position;
};

function positionToString(position: Position) {
  return `(${position.x},${position.y})`;
}

function printBoard(list: Map<string, string>, data: string[][]) {
  console.log(
    data
      .map((line, y) =>
        line
          .map(
            (tile, x) =>
              list.get(`(${x},${y})`)?.padStart(3, '0') ??
              tile.padStart(3, ' '),
          )
          .join(' '),
      )
      .join('\n'),
  );
}
