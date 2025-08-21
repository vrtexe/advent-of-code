type Direction = (typeof Direction)[keyof typeof Direction];
const Direction = {
  North: '^',
  South: 'v',
  West: '<',
  East: '>',
} as const;

const DirectionMap: Record<Direction | Slope, Position> = {
  [Direction.East]: { x: 1, y: 0 },
  [Direction.West]: { x: -1, y: 0 },
  [Direction.North]: { x: 0, y: -1 },
  [Direction.South]: { x: 0, y: 1 },
} as const;

const DirectionValues: Direction[] = Object.values(Direction);

type Slope = (typeof Slope)[keyof typeof Slope];
const Slope = {
  Up: '^',
  Down: 'v',
  Left: '<',
  Right: '>',
} as const;

function isSlope(value: string): value is Slope {
  return Slopes.includes(value);
}

const Slopes: string[] = Object.values(Slope);

function main() {
  part2_1();
}

export async function part1() {
  const data = await parse();

  const startPosition = findStartPosition(data);
  const endPosition = findEndPosition(data);

  const longestPath = findLongestPath<Point>(
    {
      position: startPosition,
      weight: 1,
      path: new Set([positionToString(startPosition)]),
    },
    (point): Point[] => {
      const [positions, skipped] = getNextTiles(
        point.position,
        point.path,
        data,
      );

      return positions.map(
        (position) =>
          <Point>{
            position,
            path: new Set([
              ...point.path,
              ...(skipped.get(positionToString(position)) ?? []),
              positionToString(position),
            ]),
            weight: point.path.size,
          },
      );
    },
    (point) => `${positionToString(point.position)}`,
  );

  const destinations = longestPath
    .values()
    .filter(
      (s) =>
        positionToString(s.vertex.position) === positionToString(endPosition),
    )
    .toArray()
    .toSorted((left, right) => right.vertex.path.size - left.vertex.path.size);

  const maxDestination = destinations[0];

  const result = maxDestination.vertex.path.size - 1;

  console.log('Result:', result);
}

type NodeGroup = {
  id: string;
  weight: number;
  group: Position[];
  edges: string[];
  path: string;
};

function buildGraph(
  data: string[][],
): [Map<string, string>, Map<string, NodeGroup>] {
  const startPosition = findStartPosition(data);
  const positionMap: Position[][] = data.map((line, y) =>
    line.map((_, x) => ({ x, y })),
  );

  const passed = new Set<Position>();
  const nextPositions = [positionMap[startPosition.y][startPosition.x]];

  const groupedPositions: Position[][] = [];
  let i = 0;

  const positionMapping = new Map<string, string>();
  const mapping = new Map<string, NodeGroup>();

  while (nextPositions.length) {
    const position = nextPositions.shift()!;

    const group: Position[] = [];
    const adjacent: Position[] = [position];

    while (adjacent.length === 1) {
      const currentAdjacent = adjacent.pop()!;

      if (!passed.has(currentAdjacent)) {
        passed.add(currentAdjacent);
      } else continue;

      group.push(currentAdjacent);

      for (const np of getNextPositionsNew(
        currentAdjacent,
        data,
        positionMap,
      )) {
        if (passed.has(np)) continue;
        adjacent.push(np);
      }
    }

    if (group.length) {
      const start = positionToString(group.at(0)!);
      const end = positionToString(group.at(-1)!);
      const key = `${start}-${end}`;

      for (const p of group) {
        positionMapping.set(positionToString(p), key);
      }

      mapping.set(key, {
        id: key,
        group: group,
        weight: group.length,
        edges: adjacent.map(positionToString),
        path: '',
      });
    }

    nextPositions.push(...adjacent);

    groupedPositions.push(group);

    // if (i > 7) {
    // console.log(position);
    // console.log(group);
    // }

    // if (i > 8) {
    //   break;
    // }

    ++i;
  }

  // for (const p of groupedPositions) {
  //   printBoard(new Set(p.map((s) => positionToString(s))), data);
  // }

  // printBoard(
  //   new Set(groupedPositions.flatMap((s) => s).map((s) => positionToString(s))),
  //   data,
  // );

  // const result = new Map<string, { group: Position[]; edges: string[] }>();
  // const edges = [positionToString(startPosition)];
  // while (edges.length) {
  //   const vertex = edges.pop()!;
  //   const current = mapping.get(vertex)!;

  //   result.set(vertex, current);
  //   edges.push(...current.edges);
  // }

  // console.log(
  //   data.flatMap((s) => s).filter((s) => ['^', 'v', '<', '>', '.'].includes(s))
  //     .length,
  // );

  for (const [, m] of mapping) {
    m.edges = m.edges.map((s) => positionMapping.get(s)!);
  }

  return [positionMapping, mapping];
}

export async function part2() {
  const data = await parse();

  const startPosition = findStartPosition(data);
  const endPosition = findEndPosition(data);

  const [positions, graph] = buildGraph(data);

  const startPositionGroup = positions.get(positionToString(startPosition))!;
  const startingGroup = graph.get(startPositionGroup)!;

  const distances = getLongestPath<Point>(
    {
      position: startPosition,
      weight: 1,
      path: new Set([positionToString(startPosition)]),
    },
    (point): Point[] => {
      const [positions, skipped] = getNextTiles(
        point.position,
        point.path,
        data,
        false,
      );

      return positions.map(
        (position) =>
          <Point>{
            position,
            path: new Set([
              ...point.path,
              ...(skipped.get(positionToString(position)) ?? []),
              positionToString(position),
            ]),
            weight: 1,
          },
      );
    },
    (p) => positionToString(p.position),

    // {
    //   id: positionToString(startPosition),
    //   weight: 0,
    //   group: [startPosition],
    //   edges: [startPositionGroup],
    //   path: positionToString(startPosition),
    // },
    // (point) => {
    //   return point.edges.map((e) => ({
    //     ...graph.get(e)!,
    //     path: `${point.path}|${e}`,
    //   }));
    // },

    // (point) => `${point.id}|${point.path}`,
    // (point) => `${extractStringFromMap(point.position, stringMap)}<-[${point.path}]`,
  );

  // const [destination] = distances
  //   .values()
  //   .filter((s) => !s.vertex.edges.length)
  //   .toArray()
  //   .toSorted((left, right) => right.weight - left.weight);

  // destination.vertex.path
  //   .split('|')
  //   .slice(1)
  //   .map((v) => graph.get(v)!)
  //   .map((v) => console.log(v.group.length));

  console.log(distances);
  // console.log(destination.weight - 1, destination.vertex.path);

  // console.log(
  //   destination,
  //   distances.get(`${destination.through!.id}|${destination.through!.path}`),
  // );
  // const destinations = longestPath
  //   .values()
  //   .filter(
  //     (s) =>
  //       positionToString(s.vertex.position) === positionToString(endPosition),
  //   )
  //   .toArray();
  // .toSorted(
  //   (left, right) =>
  //     right.vertex.path.split('-').length -
  //     left.vertex.path.split('-').length,
  // );

  // const maxDestination = destinations[0];

  // const s: Set<string> = new Set();

  // let d: Distance<Point> | undefined = maxDestination;
  // let i = 0;
  // while (d) {
  //   s.add(positionToString(d.vertex.position));
  //   d = longestPath.get(positionToString(d.through!.position!));
  //   if (i > 10000) {
  //     break;
  //   }
  //   i++;
  // }

  // printBoard(s, data);

  // > 4826
  // > 6318
  // console.log(
  //   'Result:',
  // destinations,
  // maxDestination.vertex.path.split('-').length - 1,
  // maxDestination.weight,
  // maxDestination,
  // );
}

export async function part2_1() {
  const data = await parse();

  const startPosition = findStartPosition(data);
  const endPosition = findEndPosition(data);

  // const graph = new Map<string, [Position, Set<string>, string[]]>();
  // const mapping = ""
  const state: Position[] = [startPosition];
  const passed = new Set<string>();

  const groups: Position[][] = [];

  while (state.length) {
    const position = state.shift()!;
    const positionKey = positionToString(position);

    if (!passed.has(positionKey)) {
      passed.add(positionKey);
    } else continue;

    const nextPositions: Position[] = [];
    const group: Position[] = [position];

    for (const nextPosition of findNextPositions(position, data)) {
      if (passed.has(positionToString(nextPosition))) continue;
      nextPositions.push(nextPosition);
    }

    while (nextPositions.length === 1) {
      const currentPosition = nextPositions.pop()!;

      group.push(currentPosition);
      passed.add(positionToString(currentPosition));

      for (const nextPosition of findNextPositions(currentPosition, data)) {
        if (passed.has(positionToString(nextPosition))) continue;
        nextPositions.push(nextPosition);
      }
    }

    if (group.length) {
      groups.push(group);
    }

    state.push(...nextPositions);
  }

  const groupIdOf = (g: Position[]) =>
    `${positionToString(g.at(0)!)},${positionToString(g.at(-1)!)}`;

  const tempMap = new Map<string, string>(
    groups.flatMap((g) => [
      [positionToString(g.at(0)!), groupIdOf(g)],
      [positionToString(g.at(-1)!), groupIdOf(g)],
    ]),
  );

  const graph = new Map<
    string,
    { id: string; prev: string[]; next: string[]; weight: number; path: string }
  >();

  for (const group of groups) {
    const groupId = groupIdOf(group);
    const [first, last] = [group.at(0)!, group.at(-1)!];

    const prev: string[] = [];
    const next: string[] = [];
    for (const previousPosition of findNextPositions(first, data)) {
      const prevKey = positionToString(previousPosition);
      const prevGroupId = tempMap.get(prevKey)!;
      if (!prevGroupId || groupId === prevGroupId) continue;
      prev.push(prevGroupId);
    }

    for (const nextPosition of findNextPositions(last, data)) {
      const nextKey = positionToString(nextPosition);
      const nextGroupId = tempMap.get(nextKey)!;
      if (!nextGroupId || groupId === nextGroupId) continue;
      next.push(nextGroupId);
    }

    graph.set(groupId, {
      id: groupId,
      prev,
      next,
      weight: group.length,
      path: '',
    });
  }

  graph.set(positionToString(startPosition), {
    id: positionToString(startPosition),
    next: ['(0,1),(5,3)'],
    prev: [],
    weight: 1,
    path: '',
  });
  graph.get('(0,1),(5,3)')!.weight--;
  graph.get('(0,1),(5,3)')!.prev = [positionToString(startPosition)];

  type VertexValue = Vertex & {
    value: string;
  };

  graph.get('(19,14),(19,19)')!.next = ['(20,19),(22,21)'];

  // console.log(graph);
  const startKey = positionToString(startPosition);
  const s = getLongestPath(
    {
      ...graph.get(startKey)!,
      path: startKey,
    },
    (v) => {
      // const nextPath = `${v.path}|${v.id}`;
      return (
        v.next
          .filter((s) => !v.path.includes(s))
          .map((g) => ({
            ...graph.get(g)!,
            path: `${v.path}|${g}`,
          })) ?? []
      );
    },
    (g) => `${g.id}`,
  );

  console.log('here', s.get('(20,19),(22,21)'));

  // printBoard(new Set(groups.flatMap((s) => s).map(positionToString)), data);

  // const stateMap = new Map<string, [number, number]>([
  //   [positionToString(startPosition), [0, 0]],
  // ]);
  // const state: [Position[], Set<string>, number][] = [
  //   [[startPosition], new Set<string>([positionToString(startPosition)]), 0],
  // ];

  // const results: [Position[], Set<string>, number] = [...state[0]];
  // let i = 0;

  // while (state.length) {
  //   const [positions, passed, count] = state.pop()!;
  //   const lastPosition = positions.at(-1)!;
  //   stateMap.delete(positionToString(lastPosition));

  //   if (endPosition.x === lastPosition.x && endPosition.y === lastPosition.y) {
  //     if (count > results[2]) {
  //       (results[0] = positions), (results[1] = passed), (results[2] = count);
  //       console.log(results[2], state.length);
  //     }
  //     continue;
  //   }

  //   const nextPositions: Position[] = [];

  //   for (const nextPosition of findNextPositions(lastPosition, data)) {
  //     if (passed.has(positionToString(nextPosition))) continue;
  //     nextPositions.push(nextPosition);
  //   }

  //   if (nextPositions.length === 1) {
  //     const [nextPosition] = nextPositions;
  //     const nextPositionKey = positionToString(nextPosition);
  //     const [pendingPositionCount, index] = stateMap.get(nextPositionKey) ?? [];
  //     // positions.push(nextPosition);
  //     // if (
  //     //   pendingPositionCount &&
  //     //   index != undefined &&
  //     //   pendingPositionCount < count + 1
  //     // ) {
  //     //   console.log(
  //     //     nextPositionKey == positionToString(state[index][0].at(-1)!),
  //     //   );
  //     //   state[index][1] = new Set([
  //     //     ...state[index][1],
  //     //     ...passed.add(nextPositionKey),
  //     //   ]);
  //     //   state[index][2] = count + 1;
  //     //   // pendingPositionCount
  //     // } else {
  //     passed.add(nextPositionKey);
  //     stateMap.set(nextPositionKey, [count + 1, state.length]);
  //     state.push([[nextPosition], passed, count + 1]);
  //     // }
  //   } else {
  //     // console.log('here2', nextPositions);

  //     for (const nextPosition of nextPositions) {
  //       const nextPositionState = [nextPosition];
  //       const nextPassedState = new Set(passed);
  //       const nextPositionKey = positionToString(nextPosition);
  //       const [pendingPositionCount, index] =
  //         stateMap.get(nextPositionKey) ?? [];

  //       if (
  //         index != undefined &&
  //         pendingPositionCount &&
  //         pendingPositionCount < count + 1
  //       ) {
  //         // state[index][1] = nextPassedState;
  //         state[index][2] = count + 1;
  //         // pendingPositionCount
  //       } else {
  //         // passed.add(nextPositionKey);
  //         // state.push([[nextPosition], passed, count + 1]);
  //         nextPassedState.add(nextPositionKey);
  //         stateMap.set(nextPositionKey, [count + 1, state.length]);
  //         state.push([nextPositionState, nextPassedState, count + 1]);
  //       }
  //     }
  //   }

  //   // console.log(state.length, count, i);

  //   if (i >= 16) {
  //     // console.log(state);
  //     // break;
  //   }

  //   ++i;
  // }
  // console.log(results, stateMap);
}

function extractStringFromMap(position: Position, stringMap: string[][]) {
  return stringMap[position.y]?.[position.x];
}

// deno-lint-ignore no-unused-vars
function printBoard(destination: Set<string>, board: string[][]) {
  let text = '';

  for (let y = 0; y < board.length; y++) {
    for (let x = 0; x < board[0].length; x++) {
      if (destination.has(positionToString({ y, x }))) {
        text += 'O';
        continue;
      }
      text += board[y][x];
    }
    text += '\n';
  }

  console.log(text, '\n');
}

function getNextTiles(
  position: Position,
  invalid: Set<string> | undefined,
  data: string[][],
  slopes = true,
): [Position[], Map<string, string[]>] {
  const result: Position[] = [];
  const skipped = new Map<string, string[]>();

  for (const direction of Object.values(Direction)) {
    const nextDirection = DirectionMap[direction];
    const nextPosition = addPosition(position, nextDirection);
    const tile = data[nextPosition.y]?.[nextPosition.x];

    if (!tile || invalid?.has(positionToString(nextPosition)) || tile === '#')
      continue;

    if (slopes && isSlope(tile)) {
      if (tile === oppositeDirection(direction)) continue;
      const slopeDirection = DirectionMap[tile];

      const resultPosition = addPosition(nextPosition, slopeDirection);
      const resultPositionValue = positionToString(resultPosition);
      skipped.set(resultPositionValue, [
        ...(skipped.get(resultPositionValue) ?? []),
        positionToString(nextPosition),
      ]);

      result.push(addPosition(nextPosition, slopeDirection));

      continue;
    }

    result.push(nextPosition);
  }

  return [result, skipped];
}

function efficientlyGetNextTiles(
  position: Position,
  currentPath: string,
  data: string[][],
  positionMap: Position[][],
  stringMap: string[][],
): EfficientPoint[] {
  const nextPositions = getNextPositions(
    position,
    currentPath,
    data,
    positionMap,
    stringMap,
  );

  if (nextPositions.length > 1 || nextPositions.length < 1) {
    return nextPositions.map((nextPosition) => ({
      position: nextPosition,
      path: '',
      // path: `${currentPath}-${extractStringFromMap(nextPosition, stringMap)}`,
      weight: -1,
    }));
  }

  const mergePositions: Position[] = [];
  let positions = nextPositions;
  // let currentInvalid = currentPath;

  while (positions.length && positions.length === 1) {
    const nextPosition = positions.pop()!;

    mergePositions.push(nextPosition);

    // currentInvalid = `${currentInvalid}-${extractStringFromMap(
    //   nextPosition,
    //   stringMap,
    // )}`;
    positions = getNextPositions(
      nextPosition,
      '',
      data,
      positionMap,
      stringMap,
    );
  }

  return [
    {
      position: mergePositions.at(-1)!,
      path: '',
      weight: mergePositions.reduce(count, 0),
    },
  ];
}

function count(prev: number) {
  return prev + 1;
}

function getNextPositionsNew(
  position: Position,
  data: string[][],
  positionMap: Position[][],
): Position[] {
  const result: Position[] = [];

  for (const direction of DirectionValues) {
    const nextDirection = DirectionMap[direction];
    const nextPosition = extractAddedPosition(
      position,
      nextDirection,
      positionMap,
    );
    const tile = data[nextPosition?.y]?.[nextPosition?.x];

    if (!tile || tile === '#') continue;

    result.push(nextPosition);
  }

  return result;
}

function findNextPositions(position: Position, data: string[][]): Position[] {
  const result: Position[] = [];

  for (const direction of DirectionValues) {
    const nextDirection = DirectionMap[direction];
    const nextPosition = addPosition(position, nextDirection);
    const tile = data[nextPosition?.y]?.[nextPosition?.x];

    if (!tile || tile === '#') continue;

    result.push(nextPosition);
  }

  return result;
}

function getNextPositions(
  position: Position,
  invalid: string,
  data: string[][],
  positionMap: Position[][],
  stringMap: string[][],
): Position[] {
  const result: Position[] = [];

  for (const direction of DirectionValues) {
    const nextDirection = DirectionMap[direction];
    const nextPosition = extractAddedPosition(
      position,
      nextDirection,
      positionMap,
    );
    const tile = data[nextPosition?.y]?.[nextPosition?.x];

    if (
      !tile ||
      invalid.includes(extractStringFromMap(nextPosition, stringMap)) ||
      tile === '#'
    )
      continue;

    result.push(nextPosition);
  }

  return result;
}

function oppositeDirection(direction: Direction): Direction {
  if (direction === Direction.East) return Direction.West;
  if (direction === Direction.West) return Direction.East;
  if (direction === Direction.South) return Direction.North;
  if (direction === Direction.North) return Direction.South;

  throw new Error('Invalid Direction');
}

function addPosition(left: Position, right: Position): Position {
  return { x: left.x + right.x, y: left.y + right.y };
}

function extractAddedPosition(
  left: Position,
  right: Position,
  positionGrid: Position[][],
) {
  return positionGrid[left.y + right.y]?.[left.x + right.x];
}

// deno-lint-ignore no-unused-vars
function positionFromString(value: string): Position {
  const [y, x] = value.slice(1, -1).split(',');
  return { y: parseInt(y), x: parseInt(x) };
}

function positionToString(position: Position) {
  return `(${position.y},${position.x})`;
}

// deno-lint-ignore no-unused-vars
function pathToString(path: Set<string>) {
  return path.values().toArray().join('-');
}

function findStartPosition(data: string[][]): Position {
  return { y: 0, x: data[0].indexOf('.') };
}

function findEndPosition(data: string[][]): Position {
  return { y: data.length - 1, x: data[data.length - 1].indexOf('.') };
}

async function parse() {
  const text = await Deno.readTextFile('assets/2023/task23-test.txt');
  return text.split('\n').map((line) => line.split(''));
}

if (import.meta.main) {
  main();
}

type Edge<V> = {
  from: V;
  to: V;
};

type Vertex = {
  weight: number;
};

type Position = {
  x: number;
  y: number;
};

type Point = Vertex & {
  position: Position;
  path: Set<string>;
};

type EfficientPoint = Vertex & {
  position: Position;
  path: string;
};

// deno-lint-ignore no-unused-vars
function getLongestPath<V extends Vertex>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  const distances: Map<string, Distance<V>> = new Map([
    [keyOf(vertex), { vertex, weight: 0 }],
  ]);

  const vertices: V[] = [];
  const nextVertices = [vertex];
  while (nextVertices.length) {
    const v = nextVertices.pop()!;
    vertices.push(...topologicalSort(v, getAdjacent, keyOf));
    nextVertices.push(...getAdjacent(v));
  }

  // vertices.forEach((s) => {
  // (s as unknown as { path: string }).path = (
  // s as unknown as { path: string }
  // ).path.split('|')[0];
  // });
  // for (const vertex of vertex) {

  // }

  while (vertices.length) {
    const current = vertices.shift()!;
    const currentDistance = distances.get(keyOf(current));

    if (currentDistance) {
      for (const adjacent of getAdjacent(current)) {
        const adjacentKey = keyOf(adjacent);
        const adjacentDistance = distances.get(adjacentKey);
        const nextDistance = currentDistance.weight + adjacent.weight;
        if (!adjacentDistance || adjacentDistance.weight < nextDistance) {
          distances.set(adjacentKey, {
            vertex: adjacent,
            through: current,
            weight: nextDistance,
          });
        }
      }
    }
  }

  return distances;
}

function topologicalSort<V extends string | number | object>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  const visited: Set<string> = new Set(keyOf(vertex));
  const state: V[] = [vertex];

  const result: V[] = [];

  while (state.length) {
    const current = state.pop()!;

    result.push(current);

    for (const adjacent of getAdjacent(current)) {
      const adjacentKey = keyOf(adjacent);
      if (!visited.has(adjacentKey)) {
        visited.add(adjacentKey);
        state.push(adjacent);
      }
    }
  }

  return result;
}

function findLongestPath<V extends Vertex>(
  vertex: V,
  getAdjacent: (v: V) => V[],
  keyOf = (v: V) => v.toString(),
) {
  const distances: Map<string, Distance<V>> = new Map([
    [keyOf(vertex), { vertex, weight: 0 }],
  ]);

  const queue = createPriorityQueue<[number, V]>(
    [[0, vertex]],
    ([left], [right]) => right - left,
  );

  while (queue.size()) {
    const [, current] = queue.shift()!;
    const currentDistance = distances.get(keyOf(current))!;

    const nextVertices: [number, V][] = [];

    for (const next of getAdjacent(current)) {
      const nextDistance = distances.get(keyOf(next));
      const nextWeight = currentDistance.weight + next.weight;
      if (!nextDistance || nextDistance.weight < nextWeight) {
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

// deno-lint-ignore no-unused-vars
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

type Distance<V> = {
  vertex: V;
  through?: V;
  weight: number;
};

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
