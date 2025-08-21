import { customParseInt } from '../common.ts';

function main() {
  part1();
}

export async function part1() {
  const data = await parse();
  let result = 0n;

  for (const secret of data) {
    result += findEvolution(secret, 2000);
  }

  // 11742379
  // let secret = 123;
  // console.log(secret);

  console.log('Result:', Number(result));
}

function findEvolution(secret: number, count: number) {
  let currentValue = BigInt(secret);

  for (let i = 0; i < count; i++) {
    currentValue = evolveSecretNumber(currentValue);
  }

  return currentValue;
}

function evolveSecretNumberOld(secret: number) {
  let s = secret;
  return (((s = ((s = s ^ (s * 64)) / 32) ^ s) * 2048) ^ s) % 16777216;
}

function evolveSecretNumber(secret: bigint) {
  let s = secret;
  return prune(
    mix((s = prune(mix((s = prune(mix(s, s * 64n))), s / 32n))), s * 2048n),
  );
}

function mix(left: bigint, right: bigint) {
  return left ^ right;
}

function prune(secret: bigint) {
  return secret % 16777216n;
  // return secret & 16777215;
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task22.txt');
  return text.split('\n').map(customParseInt);
}

if (import.meta.main) {
  main();
}
