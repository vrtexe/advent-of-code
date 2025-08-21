import { customParseInt } from '../common.ts';

type InstructionOpCode =
  (typeof InstructionOpCode)[keyof typeof InstructionOpCode];
const InstructionOpCode = {
  adv: 0,
  bxl: 1,
  bst: 2,
  jnz: 3,
  bxc: 4,
  out: 5,
  bdv: 6,
  cdv: 7,
} as const;

type Register = (typeof Register)[keyof typeof Register];
const Register = {
  A: 'A',
  B: 'B',
  C: 'C',
} as const;

type RegisterComboOperand =
  (typeof RegisterComboOperand)[keyof typeof RegisterComboOperand];
const RegisterComboOperand: Record<number, string> = {
  4: Register.A,
  5: Register.B,
  6: Register.C,
} as const;

function getComboOperatorValue(operand: number, state: State) {
  if (operand <= 3) return operand;

  return state.registers[RegisterComboOperand[operand]];
}

const InstructionHandler: Record<
  InstructionOpCode,
  (operand: number, state: State) => number | void
> = {
  [InstructionOpCode.adv]: (operand: number, state: State) => {
    const value = getComboOperatorValue(operand, state);
    const prev = state.registers[Register.A];
    state.registers[Register.A] = Math.trunc(prev / Math.pow(2, value));
    state.instructionPointer += 2;
  },
  [InstructionOpCode.bxl]: (operand: number, state: State) => {
    const previousValue = state.registers[Register.B];
    state.registers[Register.B] = previousValue ^ operand;
    state.instructionPointer += 2;
  },
  [InstructionOpCode.bst]: (operand: number, state: State) => {
    const value = getComboOperatorValue(operand, state);
    state.registers[Register.B] = value & 7;
    state.instructionPointer += 2;
  },
  [InstructionOpCode.jnz]: (operand: number, state: State) => {
    if (state.registers[Register.A] === 0) {
      state.instructionPointer += 2;
      return;
    }
    state.instructionPointer = operand;
  },
  [InstructionOpCode.bxc]: (_: number, state: State) => {
    state.registers[Register.B] =
      state.registers[Register.B] ^ state.registers[Register.C];
    state.instructionPointer += 2;
  },
  [InstructionOpCode.out]: (operand: number, state: State) => {
    const value = getComboOperatorValue(operand, state);
    state.instructionPointer += 2;
    return value & 7;
  },
  [InstructionOpCode.bdv]: (operand: number, state: State) => {
    const value = getComboOperatorValue(operand, state);
    const prev = state.registers[Register.A];
    state.registers[Register.B] = Math.trunc(prev / Math.pow(2, value));
    state.instructionPointer += 2;
  },
  [InstructionOpCode.cdv]: (operand: number, state: State) => {
    const value = getComboOperatorValue(operand, state);
    const prev = state.registers[Register.A];
    state.registers[Register.C] = Math.trunc(prev / Math.pow(2, value));
    state.instructionPointer += 2;
  },
};

function main() {
  part2();
}

export async function part1() {
  const state = await parse();

  const outputs = runProgram(state);
  const result = outputs.join(',');

  console.log('Result:', result);
}

function runProgram(state: State) {
  const outputs: number[] = [];

  while (state.program[state.instructionPointer] != undefined) {
    const opcode: InstructionOpCode = state.program[
      state.instructionPointer
    ] as InstructionOpCode;
    const operand: number = state.program[state.instructionPointer + 1];

    const output = InstructionHandler[opcode](operand, state);
    if (output != undefined) {
      outputs.push(output);
    }
  }

  return outputs;
}

export async function part2() {
  const data = await parse();

  const list: [number, number][] = [
    [0o1 * Math.pow(0o10, data.program.length - 1), data.program.length - 1],
  ];
  let result = 0;

  label: while (list.length) {
    const [sequence, placement] = list.shift()!;

    const nextPlacements: [number, number][] = [];

    for (let i = 0o0; i <= 0o7; i += 0o1) {
      const state = JSON.parse(JSON.stringify(data)) as State;
      const value = sequence + Math.pow(0o10, placement) * i;
      state.registers[Register.A] = value;

      const outputs = runProgram(state);

      if (outputs.join(',') === state.program.join(',')) {
        result = value;
        break label;
      }

      if (
        outputs.length === state.program.length &&
        outputs[placement] === state.program[placement]
      ) {
        nextPlacements.push([value, placement - 1]);
      }
    }

    list.push(...nextPlacements);
  }

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task17.txt');
  const [registers, program] = text.split('\n\n');
  return <State>{
    instructionPointer: 0,
    registers: Object.fromEntries(
      registers
        .split('\n')
        .map((register) => register.split(': '))
        .map(([name, value]) => [name.split(/\s+/)[1], parseInt(value)]),
    ),
    program: program.split(': ')[1].split(',').map(customParseInt),
  };
}

type State = {
  registers: Record<string, number>;
  program: number[];
  instructionPointer: number;
};

if (import.meta.main) {
  main();
}
