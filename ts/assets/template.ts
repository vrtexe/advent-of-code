function main() {
  part1();
}

export async function part1() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();
  const result = data.length;

  console.log('Result:', result);
}

async function parse() {
  const text = await Deno.readTextFile('assets/{{dataFile}}');
  return text;
}

if (import.meta.main) {
  main();
}
