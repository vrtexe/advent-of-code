function main() {
  part1();
  part2();
}

export async function part1() {
  const data = await parse();

  const [rows, columns] = [data.length, data[0].length];
  const [graph, blocks] = buildGraph(data, buildNeighborExtractor());

  const result = findShortestPath(graph, blocks, [rows, columns], data);

  console.log('Result:', result.weight);
}

export async function part2() {
  const data = await parse();

  const [rows, columns] = [data.length, data[0].length];
  const [graph, blocks] = buildGraph(
    data,
    buildNeighborExtractor({ min: 4, max: 10 }),
  );

  const result = findShortestPath(graph, blocks, [rows, columns], data);

  console.log('Result:', result.weight);
}

function vertexKey(vertex: Block) {
  return `(${vertex.position.x},${vertex.position.y},${vertex.direction})`;
}

function buildCityBlocks(direction: Direction, city: string[][]) {
  return city.map((line, y) =>
    line.map(
      (block, x) =>
        <Block>{
          position: { x, y },
          weight: parseInt(block),
          direction: direction,
        },
    ),
  );
}

function buildGraph(
  city: string[][],
  neighborExtractor: NeighborExtractor = extractNeighboringBlocks,
): [Graph<Block, Edge<Block>>, Record<Direction, Block[][]>] {
  const horizontalBlocks = buildCityBlocks(Direction.Horizontal, city);
  const verticalBlocks = buildCityBlocks(Direction.Vertical, city);

  const blocks: Record<Direction, Block[][]> = {
    [Direction.Horizontal]: horizontalBlocks,
    [Direction.Vertical]: verticalBlocks,
  };

  const graph = new Graph(
    new Set(
      Object.values(blocks)
        .flatMap((s) => s)
        .flat(),
    ),
    new Map(),
    (v: Block) => neighborExtractor(v, blocks),
    vertexKey,
  );

  return [graph, blocks];
}

function findShortestPath(
  graph: Graph<Block, Edge<Block>>,
  blocks: Record<Direction, Block[][]>,
  [rows, columns]: [number, number],
  city: string[][],
  print = false,
) {
  const distances = graph.findShortestPath({
    ...blocks[Direction.Horizontal][0][0],
    weight: 0,
    direction: Direction.Horizontal,
  });

  const destinations: Distance<Block>[] = [
    blocks[Direction.Horizontal][rows - 1][columns - 1],
    blocks[Direction.Vertical][rows - 1][columns - 1],
  ]
    .map((block) => distances.get(vertexKey(block)))
    .filter((s) => s != undefined)
    .toSorted((left, right) => left.weight - right.weight);

  print && applyAndPrint(destinations[0], distances, city);

  return destinations[0];
}

function buildNeighborExtractor(
  constraint: Constraint = { min: 1, max: 3 },
): NeighborExtractor {
  return (block, blocks) => {
    return extractNeighboringBlocks(block, blocks, constraint);
  };
}

function extractNeighboringBlocks(
  block: Block,
  blocks: Record<Direction, Block[][]>,
  constraint: Constraint = { min: 1, max: 3 },
) {
  const position = block.position;

  return extractNextNeighbors(
    block,
    constraint,
    createNeighborResolver(position, block.direction!, blocks),
  );
}

function createNeighborResolver(
  position: Position,
  direction: Direction,
  blocks: Record<Direction, Block[][]>,
): (increment: number) => Block {
  const flippedDirection = nextDirection(direction);

  switch (direction) {
    case Direction.Horizontal:
      return (i) => blocks[flippedDirection][position.y]?.[position.x + i];
    case Direction.Vertical:
      return (i: number) =>
        blocks[flippedDirection][position.y + i]?.[position.x];
  }
}

function extractNextNeighbors(
  block: Block,
  { min, max }: Constraint,
  next: (increment: number) => Block,
) {
  const direction = nextDirection(block.direction!);
  return [
    ...calculateNegativeWeights(-min, -max, direction, next),
    ...calculatePositiveWeights(min, max, direction, next),
  ];
}

function calculateNegativeWeights(
  start: number,
  end: number,
  direction: Direction,
  next: (increment: number) => Block,
) {
  let weight = 0;
  const result: Block[] = [];

  for (let i = -1; i > start; --i) {
    const nextBlock = next(i);
    if (!nextBlock) break;
    weight += nextBlock.weight;
  }

  for (let i = start; i >= end; --i) {
    const nextBlock = next(i);
    if (!nextBlock) break;

    weight += nextBlock.weight;
    result.push({
      ...nextBlock,
      direction,
      weight,
    });
  }

  return result;
}

function calculatePositiveWeights(
  start: number,
  end: number,
  direction: Direction,
  next: (increment: number) => Block,
) {
  let weight = 0;
  const result: Block[] = [];

  for (let i = 1; i < start; ++i) {
    const nextBlock = next(i);
    if (!nextBlock) break;
    weight += nextBlock.weight;
  }

  for (let i = start; i <= end; ++i) {
    const nextBlock = next(i);
    if (!nextBlock) break;

    weight += nextBlock.weight;
    result.push({
      ...nextBlock,
      direction,
      weight,
    });
  }

  return result;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2023/task17.txt');
  return text.split('\n').map((line) => line.split(''));
}

if (import.meta.main) {
  main();
}

function deepCopy<T>(value: T): T {
  return JSON.parse(JSON.stringify(value));
}

function applyAndPrint(
  startingBlock: Distance<Block>,
  destinations: Map<string, Distance<Block>>,
  city: string[][],
) {
  const board = deepCopy(city);
  let current: Distance<Block> | undefined = startingBlock;
  const path: Distance<Block>[] = [];
  while (current) {
    path.push(current);

    if (!current.through) break;

    current = destinations.get(vertexKey(current.through));
  }

  for (let i = 1; i < path.length; i++) {
    const to = path[i - 1];
    const from = path[i];

    let position = '';
    if (from.vertex.position.x < to.vertex.position.x) {
      position = '>';
    } else if (from.vertex.position.x > to.vertex.position.x) {
      position = '<';
    } else if (from.vertex.position.y < to.vertex.position.y) {
      position = 'v';
    } else if (from.vertex.position.y > to.vertex.position.y) {
      position = '^';
    } else {
      throw new Error('Invalid positions');
    }

    for (
      let x = Math.min(from.vertex.position.x, to.vertex.position.x);
      x <= Math.max(from.vertex.position.x, to.vertex.position.x);
      x++
    ) {
      board[from.vertex.position.y][x] = position;
    }
    for (
      let y = Math.min(from.vertex.position.y, to.vertex.position.y);
      y <= Math.max(from.vertex.position.y, to.vertex.position.y);
      y++
    ) {
      board[y][from.vertex.position.x] = position;
    }
  }

  printBoard(board);
}

function printBoard(board: string[][]) {
  for (const row of board) {
    console.log(row.join(''));
  }

  console.log();
}

type Edge<V> = {
  from: V;
  to: V;
};

type Vertex = {
  weight: number;
};

type NeighborExtractor = (
  block: Block,
  blocks: Record<Direction, Block[][]>,
) => Block[];

class Graph<V extends Vertex, E extends Edge<V>> {
  private vertices: Set<V>;
  private edges: Map<V, E[]>;

  private customAdjacentProvider: ((vertex: V) => V[]) | undefined;
  private customVertexKey: ((vertex: V) => string) | undefined;

  constructor(
    vertices: Set<V>,
    edges: Map<V, E[]>,
    customAdjacentProvider?: (vertex: V) => V[],
    customVertexKey?: ((vertex: V) => string) | undefined,
  ) {
    this.vertices = vertices;
    this.edges = edges;
    this.customAdjacentProvider = customAdjacentProvider;
    this.customVertexKey = customVertexKey;
  }

  public findShortestPath(vertex: V) {
    const distances: Map<string, Distance<V>> = new Map([
      [this.keyOf(vertex), { vertex, weight: 0 }],
    ]);

    const queue = createPriorityQueue<[number, V]>(
      [[0, vertex]],
      ([left], [right]) => left - right,
    );

    while (queue.size()) {
      const [, current] = queue.shift()!;
      const currentDistance = distances.get(this.keyOf(current))!;

      const nextVertices: [number, V][] = [];

      for (const next of this.getAdjacent(current)) {
        const nextDistance = distances.get(this.keyOf(next));
        if (
          !nextDistance ||
          nextDistance.weight > currentDistance.weight + next.weight
        ) {
          const nextWeight = currentDistance.weight + next.weight;
          distances.set(this.keyOf(next), {
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

  public getAdjacent(vertex: V) {
    if (this.customAdjacentProvider) {
      return this.customAdjacentProvider(vertex);
    }

    return this.edges.get(vertex)?.map((e) => e.to) ?? [];
  }

  private keyOf(vertex: V) {
    if (this.customVertexKey) {
      return this.customVertexKey(vertex);
    }

    return JSON.stringify(vertex);
  }
}

type Comparator<T> = (left: T, right: T) => number;

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

type Position = {
  x: number;
  y: number;
};

enum Direction {
  Horizontal = 'H',
  Vertical = 'V',
}

type Constraint = {
  min: number;
  max: number;
};

function nextDirection(direction: Direction) {
  if (direction === Direction.Horizontal) return Direction.Vertical;
  if (direction === Direction.Vertical) return Direction.Horizontal;

  throw new Error('Invalid Direction');
}

type Block = Vertex & {
  position: Position;
  direction?: Direction;
};
