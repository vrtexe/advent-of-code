function main() {
  part2();
}

export async function part1() {
  const [board, antennas] = await parse();
  let result = 0;

  for (const key in antennas) {
    const tunedAntennas = [...antennas[key]];

    while (tunedAntennas.length) {
      const antenna = tunedAntennas.pop()!;
      for (const nextAntenna of tunedAntennas) {
        const projections = findProjections(
          antenna.position,
          nextAntenna.position,
        );

        for (const projection of projections) {
          if (board[projection.y]?.[projection.x]?.data?.antinode === false) {
            board[projection.y][projection.x].data.antinode = true;
            result++;
          }
        }
      }
    }
  }

  printBoard(board);

  console.log('Result:', result);
}

export async function part2() {
  const [board, antennas] = await parse();
  let result = 0;

  const boardBoundary: BoardBoundary = [board.length, board[0].length];

  for (const key in antennas) {
    const tunedAntennas = [...antennas[key]];

    while (tunedAntennas.length) {
      const antenna = tunedAntennas.pop()!;
      for (const nextAntenna of tunedAntennas) {
        const projections = findAllProjections(
          antenna.position,
          nextAntenna.position,
          boardBoundary,
        );

        for (const projection of projections) {
          if (board[projection.y]?.[projection.x]?.data?.antinode === false) {
            board[projection.y][projection.x].data.antinode = true;
            result++;
          }
        }
      }
    }
  }

  printBoard(board, false);

  console.log('Result:', result);
}

function findProjections(
  left: Position,
  right: Position,
): [Position, Position] {
  const [xDifference, yDifference] = findDifference(left, right);

  return [
    {
      x: left.x - xDifference,
      y: left.y - yDifference,
    },
    {
      x: right.x + xDifference,
      y: right.y + yDifference,
    },
  ];
}

function findAllProjections(
  left: Position,
  right: Position,
  boardBoundary: BoardBoundary,
) {
  const [xDifference, yDifference] = findDifference(left, right);
  const leftProjections = interpolateProjections(
    left,
    [-xDifference, -yDifference],
    boardBoundary,
  );
  const rightProjections = interpolateProjections(
    right,
    [xDifference, yDifference],
    boardBoundary,
  );

  return [...leftProjections, ...rightProjections];
}

function interpolateProjections(
  start: Position,
  [xDifference, yDifference]: PositionDifference,
  [rows, columns]: BoardBoundary,
) {
  const result: Position[] = [start];

  while (true) {
    const projection = result.at(-1)!;

    const nextProjection = {
      x: projection.x + xDifference,
      y: projection.y + yDifference,
    };

    if (nextProjection.x >= rows || nextProjection.x < 0) break;
    if (nextProjection.y >= columns || nextProjection.y < 0) break;

    result.push(nextProjection);
  }

  return result;
}

function findDifference(left: Position, right: Position): PositionDifference {
  return [right.x - left.x, right.y - left.y];
}

async function parse(): Promise<[NodeBoard, AntennaMap]> {
  const text = await Deno.readTextFile('src/data/task8.txt');
  const antennas: Record<string, Antenna[]> = {};
  const data = text.split('\n').map((line, y) =>
    line.split('').map((node, x) => {
      if (node === '.') {
        return <Node>{
          type: NodeType.Empty,
          data: { antinode: false },
        };
      }

      const antenna: Antenna = {
        frequency: node,
        position: { x, y },
        antinode: false,
      };

      (antennas[node] = antennas[node] || []).push(antenna);

      return <Node>{
        type: NodeType.Antenna,
        data: antenna,
      };
    }),
  );
  return [data, antennas];
}

function printBoard(board: NodeBoard, printAntennaAntiNodes = false) {
  for (const row of board) {
    console.log(
      row
        .map((node) => {
          switch (node.type) {
            case NodeType.Empty:
              return node.data.antinode ? '#' : '.';
            case NodeType.Antenna:
              if (printAntennaAntiNodes && node.data.antinode) return '#';
              return node.data.frequency;
          }
        })
        .join(''),
    );
  }

  console.log();
}

// deno-lint-ignore no-unused-vars
function generateBoard(rows: number, columns: number): string[][] {
  return new Array(rows).fill([]).map(() => new Array(columns).fill('.'));
}

if (import.meta.main) {
  main();
}

enum NodeType {
  Empty = 'EMPTY',
  Antenna = 'ANTENNA',
}

type BaseNode = {
  antinode: boolean;
};

type Node =
  | {
      type: NodeType.Antenna;
      data: BaseNode & Antenna;
    }
  | {
      type: NodeType.Empty;
      data: BaseNode;
    };

type BoardBoundary = [number, number];
type PositionDifference = [number, number];
type NodeBoard = Node[][];
type AntennaMap = Record<string, Antenna[]>;

type Position = {
  x: number;
  y: number;
};

type Antenna = {
  frequency: string;
  position: Position;
  antinode: boolean;
};
