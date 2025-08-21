function main() {
  part1();
  part2();
}

export async function part1() {
  const data = await parse();

  const cycle = 100;

  // const [rows, cols] = [7, 11];
  const [rows, cols] = [103, 101];

  const [middleRow, middleCol] = [(rows - 1) / 2, (cols - 1) / 2];

  let [q1, q2, q3, q4] = [0, 0, 0, 0];

  for (const robot of data) {
    const robotPosition = findPositionInCycle(cycle, robot, [rows, cols]);

    if (robotPosition.y < middleRow && robotPosition.x < middleCol) ++q1;
    if (robotPosition.y > middleRow && robotPosition.x < middleCol) ++q2;
    if (robotPosition.y > middleRow && robotPosition.x > middleCol) ++q3;
    if (robotPosition.y < middleRow && robotPosition.x > middleCol) ++q4;
  }

  const result = q1 * q2 * q3 * q4;

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const [rows, cols] = [103, 101];

  const positions: [Position, Position][][] = [
    data.map((s) => [s.start, s.velocity]),
  ];

  let result = 0;

  while (true) {
    const lastState = positions.at(-1)!;
    const nextPositions: [Position, Position][] = [];
    for (const [position, velocity] of lastState) {
      const next = {
        x: teleport(position.x + velocity.x, cols),
        y: teleport(position.y + velocity.y, rows),
      };

      nextPositions.push([next, velocity]);
    }

    ++result;

    positions.push(nextPositions);

    if (isResultFrame(nextPositions, [rows, cols])) break;
  }

  console.log('Result:', result);
}

function isResultFrame(
  positions: [Position, Position][],
  [rows, cols]: [number, number],
) {
  const board: string[][] = new Array(rows)
    .fill([])
    .map((): string[] => new Array(cols).fill('.'));

  for (const [position] of positions) {
    board[position.y][position.x] = 'O';
  }

  const boardText = board.map((line) => line.join('')).join('\n');
  return boardText.includes('OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO');
}

// deno-lint-ignore no-unused-vars
function printBoard(
  positions: [Position, Position][],
  [rows, cols]: [number, number],
) {
  const board: string[][] = new Array(rows)
    .fill([])
    .map((): string[] => new Array(cols).fill('.'));

  for (const [position] of positions) {
    board[position.y][position.x] = 'O';
  }

  const boardText = board.map((line) => line.join('')).join('\n');

  console.log(boardText);
}

function findPositionInCycle(
  cycle: 100,
  robot: Robot,
  [rows, cols]: [number, number],
) {
  let position = robot.start;
  for (let i = 0; i < cycle; i++) {
    const next = {
      x: teleport(position.x + robot.velocity.x, cols),
      y: teleport(position.y + robot.velocity.y, rows),
    };
    position = next;
  }

  return position;
}

// deno-lint-ignore no-unused-vars
function fillInPositions(robot: Robot, [rows, cols]: [number, number]) {
  const positions = [robot.start];
  const pastPositions: string[] = [positionToString(robot.start)];

  let repeatStart = 0;

  while (true) {
    const position = positions.at(-1)!;
    const next = {
      x: teleport(position.x + robot.velocity.x, cols),
      y: teleport(position.y + robot.velocity.y, rows),
    };

    const positionString = positionToString(next);
    repeatStart = pastPositions.indexOf(positionString);

    if (repeatStart >= 0) break;

    pastPositions.push(positionString);
    positions.push(next);
  }

  robot.positions = positions;
  robot.cycleStart = positions.length - repeatStart;

  return robot;
}

function positionToString(position: Position) {
  return `(${position.x},${position.y})`;
}

function teleport(value: number, maxValue: number) {
  const newValue = value % maxValue;
  return newValue >= 0 ? newValue : maxValue + newValue;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task14.txt');
  return text.split('\n').map(parseRobot);
}

function parseRobot(line: string): Robot {
  const [, px, py, vx, vy] = /^p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)/.exec(line)!;

  return {
    start: { x: parseInt(px), y: parseInt(py) },
    velocity: { x: parseInt(vx), y: parseInt(vy) },
    positions: [],
    cycleStart: 0,
  };
}

type Position = {
  y: number;
  x: number;
};

type Robot = {
  start: Position;
  velocity: Position;

  positions: Position[];
  cycleStart: number;
};

if (import.meta.main) {
  main();
}
