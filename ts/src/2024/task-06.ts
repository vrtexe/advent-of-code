type TileType = (typeof TileType)[keyof typeof TileType];
const TileType = {
  Obstruction: '#',
  Guard: '^',
  Empty: '.',
  PlacedObstruction: 'O',
} as const;

const Obstructions: TileType[] = [
  TileType.Obstruction,
  TileType.PlacedObstruction,
];

type Direction = (typeof Direction)[keyof typeof Direction];
const Direction = {
  Up: '^',
  Down: 'v',
  Left: '<',
  Right: '>',
} as const;

function main() {
  part2();
}

export async function part1() {
  const [board, startingGuardTile] = await parse();

  const path = guardWalkingPath(board, startingGuardTile);

  console.log('Result:', path?.size);
}

export async function part2() {
  const [board, startingGuardTile] = await parse();
  let result = 0;

  const clonedBoard = deepCopy(board);
  const path = guardWalkingPath(
    clonedBoard,
    clonedBoard[startingGuardTile.position.y][startingGuardTile.position.x],
  )!;

  let guardTile = startingGuardTile;
  for (const tile of path) {
    if (isSamePosition(tile, startingGuardTile)) continue;
    const { position } = tile;
    const obstructedBoard = placeObstruction(position, deepCopy(board));

    const obstructedPath = guardWalkingPath(
      obstructedBoard,
      obstructedBoard[guardTile.position.y][guardTile.position.x],
    );

    guardTile.guard!.direction = interpolateDirection(
      guardTile.position,
      position,
    );
    guardTile = updateGuardPosition(position, guardTile, board)!;

    result += !obstructedPath ? 1 : 0;
  }

  console.log('Result:', result);
}

function interpolateDirection(left: Position, right: Position) {
  if (left.x < right.x) {
    return Direction.Right;
  } else if (left.x > right.x) {
    return Direction.Left;
  }

  if (left.y < right.y) {
    return Direction.Down;
  } else {
    return Direction.Up;
  }
}

function isSamePosition(left: Tile, right: Tile) {
  return (
    left.position.x === right.position.x && left.position.y === right.position.y
  );
}

function placeObstruction(position: Position, board: Board) {
  board[position.y][position.x].type = TileType.PlacedObstruction;
  return board;
}

function deepCopy<T>(value: T) {
  return JSON.parse(JSON.stringify(value));
}

function guardWalkingPath(board: Board, startingGuardTile: Tile) {
  const visitedPosition: Set<Tile> = new Set([startingGuardTile]);
  let guardTile = startingGuardTile;

  while (true) {
    const nextPosition = moveGuard(guardTile);

    if (!board[nextPosition.y]?.[nextPosition.x]) break;

    const nextGuardTile = updateGuardPosition(nextPosition, guardTile, board);
    if (!nextGuardTile) return;

    guardTile = nextGuardTile;
    visitedPosition.add(guardTile);
  }

  return visitedPosition;
}

// deno-lint-ignore no-unused-vars
function printBoard(board: Board) {
  for (const row of board) {
    console.log(
      row
        .map((s) => {
          if (s.guard) {
            return s.guard.direction;
          }
          if (s.passed) {
            if (
              (s.moved.directions.includes(Direction.Up) ||
                s.moved.directions.includes(Direction.Down)) &&
              (s.moved.directions.includes(Direction.Left) ||
                s.moved.directions.includes(Direction.Right))
            )
              return '+';
            if (
              s.moved.directions.includes(Direction.Up) ||
              s.moved.directions.includes(Direction.Down)
            )
              return '|';

            if (
              s.moved.directions.includes(Direction.Left) ||
              s.moved.directions.includes(Direction.Right)
            )
              return '-';
            return 'X';
          }
          return s.type;
        })
        .join(''),
    );
  }

  console.log();
}

function updateGuardPosition(
  nextPosition: Position,
  guardPosition: Tile,
  state: Board,
): Tile | undefined {
  if (!guardPosition.guard) throw new Error('Guard cant be null');

  if (!guardPosition.moved.directions.includes(guardPosition.guard.direction)) {
    guardPosition.moved.directions.push(guardPosition.guard.direction);
  }

  const nextTile = state[nextPosition.y][nextPosition.x];

  if (Obstructions.includes(nextTile.type)) {
    guardPosition.guard.direction = nextDirection(
      guardPosition.guard.direction,
    );
    return guardPosition;
  }

  if (nextTile.passed?.directions.includes(guardPosition.guard.direction)) {
    return;
  }

  nextTile.guard = guardPosition.guard;

  if (nextTile.passed) {
    nextTile.passed.directions.push(guardPosition.guard.direction);
  } else {
    nextTile.passed = { directions: [guardPosition.guard.direction] };
  }

  guardPosition.guard = undefined;

  return nextTile;
}

const directionOrder = [
  Direction.Up,
  Direction.Right,
  Direction.Down,
  Direction.Left,
];

function nextDirection(direction: Direction): Direction {
  return directionOrder[
    (directionOrder.indexOf(direction) + 1) % directionOrder.length
  ];
}

function moveGuard(guardPosition: Tile): Position {
  const guard = guardPosition.guard;
  if (!guard) throw new Error('Guard cant be null');

  switch (guard.direction) {
    case Direction.Up:
      return {
        x: guardPosition.position.x,
        y: guardPosition.position.y - 1,
      };
    case Direction.Down:
      return {
        x: guardPosition.position.x,
        y: guardPosition.position.y + 1,
      };
    case Direction.Left:
      return {
        x: guardPosition.position.x - 1,
        y: guardPosition.position.y,
      };
    case Direction.Right:
      return {
        x: guardPosition.position.x + 1,
        y: guardPosition.position.y,
      };
  }
}

async function parse(): Promise<[Board, Tile]> {
  const text = await Deno.readTextFile('src/data/task6.txt');

  let guardPosition: Tile;
  const data = text.split('\n').map((line, y) => {
    return line.split('').map((tile, x) => {
      if (tile === TileType.Guard) {
        return (guardPosition = <Tile>{
          position: { x, y },
          type: tile,
          passed: { directions: [Direction.Up] },
          moved: { directions: [] },
          guard: {
            direction: tile,
          },
        });
      }
      return <Tile>{
        position: { x, y },
        type: tile,
        moved: { directions: [] },
      };
    });
  });

  return [data, guardPosition!];
}

if (import.meta.main) {
  main();
}

type Position = {
  x: number;
  y: number;
};

type Tile = {
  position: Position;
  type: TileType;
  passed?: {
    directions: Direction[];
  };
  moved: {
    directions: Direction[];
  };
  guard?: Guard;
};

type Guard = {
  direction: Direction;
};

type Board = Tile[][];
