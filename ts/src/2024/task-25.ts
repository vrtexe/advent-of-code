const FilledRow = '#'.repeat(5);
const EmptyRow = '.'.repeat(5);

function main() {
  part1();
}

export async function part1() {
  const data = await parse();

  const keys = data.filter((s) => s.type === SchematicType.Key);
  const locks = data.filter((s) => s.type === SchematicType.Lock);

  let result = 0;

  for (const key of keys) {
    for (const lock of locks) {
      if (testFit(key, lock)) {
        result++;
      }
    }
  }

  console.log('Result:', result);
}

function testFit(left: Schematic, right: Schematic) {
  if (left.code.length !== right.code.length) return false;
  for (let i = 0; i < left.code.length; i++) {
    if (left.code[i] + right.code[i] > 5) {
      return false;
    }
  }

  return true;
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task25.txt');
  return text.split('\n\n').map(parseSchematic);
}

function parseSchematic(value: string) {
  const lines = value.split('\n');
  const type = findSchematicType(lines);

  const code: number[] = [];

  for (let row = 0; row < lines.length; row++) {
    for (let col = 0; col < lines[row].length; col++) {
      code[col] = (code[col] ?? 0) + (lines[row][col] === '#' ? 1 : 0);
    }
  }

  return <Schematic>{
    type,
    code: code.map((s) => s - 1),
  };
}

function findSchematicType(schematic: string[]): SchematicType {
  if (
    schematic[0].includes(FilledRow) &&
    schematic[schematic.length - 1].includes(EmptyRow)
  ) {
    return SchematicType.Lock;
  }

  if (
    schematic[0].includes(EmptyRow) &&
    schematic[schematic.length - 1].includes(FilledRow)
  ) {
    return SchematicType.Key;
  }

  throw new Error('Invalid Schematic');
}

type SchematicType = (typeof SchematicType)[keyof typeof SchematicType];
const SchematicType = {
  Lock: 'LOCK',
  Key: 'KEY',
} as const;

type Schematic = {
  type: SchematicType;
  code: number[];
};

if (import.meta.main) {
  main();
}
