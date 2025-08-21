import { customParseInt, multiply, sum } from '../common.ts';

type Operation = (typeof Operation)[keyof typeof Operation];
const Operation = Object.freeze({
  MULTIPLY: 'mul',
  ENABLE: 'do',
  DISABLE: "don't",
} as const);
type Instruction = {
  operation: string;
  args: string[];
};

function main() {
  part2();
}

export async function part1() {
  const instructions = await parse();
  const result = instructions.map(execInstruction).reduce(sum);

  console.log('Result:', result);
}

function execInstruction({ operation, args }: Instruction) {
  switch (operation) {
    case Operation.MULTIPLY:
      return multiplyInstruction(args);
  }

  return 0;
}

function multiplyInstruction(args: string[]) {
  return args.map(customParseInt).reduce(multiply, 1);
}

export async function part2() {
  const instructions = await parse();

  let enabled = true;
  let result = 0;

  for (const instruction of instructions) {
    enabled = execStateInstruction(instruction) ?? enabled;
    if (!enabled) continue;

    result += execInstruction(instruction);
  }

  console.log('Result:', result);
}

function execStateInstruction(instruction: Instruction) {
  switch (instruction.operation) {
    case Operation.ENABLE:
      return true;
    case Operation.DISABLE:
      return false;
  }
}

async function parse(): Promise<Instruction[]> {
  const text = await Deno.readTextFile('src/data/task3.txt');
  return text.match(/(mul\(\d+,\d+\))|(do\(\))|(don't\(\))/g)!.map((v) => {
    const [operation, args] = v.split('(');
    return {
      operation,
      args: args
        .slice(0, -1)
        .split(',')
        .filter((s) => s),
    };
  });
}

if (import.meta.main) {
  main();
}
