import { sum } from '../common.ts';

type Movement = (typeof Movement)[keyof typeof Movement];
const Movement = {
  Up: '^',
  Down: 'v',
  Left: '<',
  Right: '>',
} as const;

type Tile = (typeof Tile)[keyof typeof Tile];
const Tile = {
  Robot: '@',
  Box: 'O',
  Wall: '#',
  Empty: '.',
  BoxLeft: '[',
  BoxRight: ']',
};

const BigBox: string[] = [Tile.BoxLeft, Tile.BoxRight] as const;

const MovementDirection: Record<Movement, Position> = {
  [Movement.Up]: { x: 0, y: -1 },
  [Movement.Down]: { x: 0, y: 1 },
  [Movement.Left]: { x: -1, y: 0 },
  [Movement.Right]: { x: 1, y: 0 },
};

function main() {
  part2();
}

export async function part1() {
  const [grid, movements] = await parse();

  let robotPosition = findRobotPosition(grid);

  for (const movement of movements) {
    const nextRobotPosition = addPosition(
      robotPosition,
      MovementDirection[movement],
    );

    if (grid[nextRobotPosition.y][nextRobotPosition.x] === Tile.Wall) {
      continue;
    }

    if (grid[nextRobotPosition.y][nextRobotPosition.x] === Tile.Empty) {
      grid[robotPosition.y][robotPosition.x] = Tile.Empty;
      grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
      robotPosition = nextRobotPosition;
      continue;
    }

    const boxesInPath: Position[] = [];
    let currentBoxPosition = nextRobotPosition;

    while (grid[currentBoxPosition.y][currentBoxPosition.x] === Tile.Box) {
      boxesInPath.push(currentBoxPosition);
      currentBoxPosition = addPosition(
        currentBoxPosition,
        MovementDirection[movement],
      );
    }

    if (grid[currentBoxPosition.y][currentBoxPosition.x] === Tile.Empty) {
      grid[robotPosition.y][robotPosition.x] = Tile.Empty;
      grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
      robotPosition = nextRobotPosition;
      for (const boxPosition of boxesInPath) {
        const newBoxPosition = addPosition(
          boxPosition,
          MovementDirection[movement],
        );
        grid[newBoxPosition.y][newBoxPosition.x] = Tile.Box;
      }
    }
  }

  const result = grid
    .flatMap((line, y) => {
      return line.map((tile, x) => {
        if (tile === Tile.Box) return 100 * y + x;
        return 0;
      });
    })
    .reduce(sum, 0);

  // console.log(grid.map((line) => line.join('')).join('\n'));
  // console.log();

  console.log('Result:', result);
}

function addPosition(left: Position, right: Position): Position {
  return { x: left.x + right.x, y: left.y + right.y };
}

function findRobotPosition(data: string[][]): Position {
  for (let i = 0; i < data.length; i++) {
    for (let j = 0; j < data[i].length; j++) {
      if (data[i][j] === Tile.Robot) {
        return { y: i, x: j };
      }
    }
  }

  throw new Error('Could not find the robot tile');
}

export async function part2() {
  const [data, movements] = await parse();

  const grid = widenGrid(data);

  let robotPosition = findRobotPosition(grid);

  for (const movement of movements) {
    const nextRobotPosition = addPosition(
      robotPosition,
      MovementDirection[movement],
    );

    if (grid[nextRobotPosition.y][nextRobotPosition.x] === Tile.Wall) {
      continue;
    }

    if (grid[nextRobotPosition.y][nextRobotPosition.x] === Tile.Empty) {
      grid[robotPosition.y][robotPosition.x] = Tile.Empty;
      grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
      robotPosition = nextRobotPosition;
      continue;
    }

    const boxesInPath: [Position, Tile][] = [];
    const currentBoxPositions: Position[] = [nextRobotPosition];
    let afterBoxTile: string | undefined;

    while (currentBoxPositions.length) {
      const currentBoxPosition = currentBoxPositions.pop()!;
      const [otherBoxPosition, otherBoxTile] = getNextBoxPart(
        currentBoxPosition,
        grid[currentBoxPosition.y][currentBoxPosition.x],
      );
      boxesInPath.push([
        currentBoxPosition,
        grid[currentBoxPosition.y][currentBoxPosition.x],
      ]);

      boxesInPath.push([otherBoxPosition, otherBoxTile]);
      if ((<Movement[]>[Movement.Left, Movement.Right]).includes(movement)) {
        const nextPos = addPosition(
          otherBoxPosition,
          MovementDirection[movement],
        );
        if (BigBox.includes(grid[nextPos.y][nextPos.x])) {
          currentBoxPositions.push(nextPos);
        } else {
          afterBoxTile = grid[nextPos.y][nextPos.x];
          break;
        }
      } else {
        const currentNextPosition = addPosition(
          currentBoxPosition,
          MovementDirection[movement],
        );
        const otherNextPosition = addPosition(
          otherBoxPosition,
          MovementDirection[movement],
        );

        if (BigBox.includes(grid[currentNextPosition.y][currentNextPosition.x]))
          currentBoxPositions.push(currentNextPosition);
        if (BigBox.includes(grid[otherNextPosition.y][otherNextPosition.x])) {
          currentBoxPositions.push(otherNextPosition);
        }
      }
    }

    if (afterBoxTile && afterBoxTile === Tile.Empty) {
      grid[robotPosition.y][robotPosition.x] = Tile.Empty;
      grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
      robotPosition = nextRobotPosition;
      for (const [boxPosition, tile] of boxesInPath) {
        const newBoxPosition = addPosition(
          boxPosition,
          MovementDirection[movement],
        );
        grid[newBoxPosition.y][newBoxPosition.x] = tile;
      }
    } else {
      if (
        boxesInPath
          .map(([p]) => {
            const nextPosition = addPosition(p, MovementDirection[movement]);
            return grid[nextPosition.y][nextPosition.x] !== Tile.Wall;
          })
          .reduce((prev, next) => prev && next, true)
      ) {
        if (movement === Movement.Up) {
          grid[robotPosition.y][robotPosition.x] = Tile.Empty;
          grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
          robotPosition = nextRobotPosition;
        }
        const data = boxesInPath.toSorted(([left], [right]) => {
          return movement === Movement.Up ? left.y - right.y : right.y - left.y;
        });

        for (const [boxPosition, tile] of data) {
          const newBoxPosition = addPosition(
            boxPosition,
            MovementDirection[movement],
          );
          grid[newBoxPosition.y][newBoxPosition.x] = tile;
          grid[boxPosition.y][boxPosition.x] = '.';
        }

        if (movement === Movement.Down) {
          grid[robotPosition.y][robotPosition.x] = Tile.Empty;
          grid[nextRobotPosition.y][nextRobotPosition.x] = Tile.Robot;
          robotPosition = nextRobotPosition;
        }
      }
    }
  }

  const result = grid
    .flatMap((line, y) => {
      return line.map((tile, x) => {
        if (tile === Tile.BoxLeft) return 100 * y + x;
        return 0;
      });
    })
    .reduce(sum, 0);

  // console.log(grid.map((line) => line.join('')).join('\n'));
  // console.log();

  console.log('Result:', result);
}

function getNextBoxPart(position: Position, tile: string): [Position, Tile] {
  if (tile === Tile.BoxLeft) {
    return [
      addPosition(position, MovementDirection[Movement.Right]),
      Tile.BoxRight,
    ];
  }

  return [
    addPosition(position, MovementDirection[Movement.Left]),
    Tile.BoxLeft,
  ];
}

function widenGrid(grid: string[][]) {
  return grid.map((line) =>
    line.flatMap((tile) => {
      switch (tile) {
        case Tile.Box:
          return [Tile.BoxLeft, Tile.BoxRight];
        case Tile.Wall:
          return [Tile.Wall, Tile.Wall];
        case Tile.Robot:
          return [Tile.Robot, Tile.Empty];
        case Tile.Empty:
          return [Tile.Empty, Tile.Empty];
      }
      throw new Error('Invalid tile');
    }),
  );
}

async function parse(): Promise<[string[][], Movement[]]> {
  const text = await Deno.readTextFile('assets/2024/task15.txt');
  const [board, movements] = text.split('\n\n');
  return [
    board.split('\n').map((b) => b.split('')),
    movements.split('\n').flatMap((line) => line.split('')) as Movement[],
  ];
}

if (import.meta.main) {
  main();
}

type Position = {
  x: number;
  y: number;
};
