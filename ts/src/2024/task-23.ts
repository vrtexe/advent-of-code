function main() {
  part2();
}

export async function part1() {
  const data = await parse();

  const edges = new Set<string>();
  for (const [k, v] of data) {
    if (k.startsWith('t')) {
      for (const left of v) {
        for (const right of v) {
          if (left === right) continue;
          if (!data.get(left)?.has(right) || !data.get(right)?.has(left))
            continue;
          edges.add([k, left, right].toSorted().join(','));
        }
      }
    }
  }

  const result = edges.size;

  console.log('Result:', result);
}

export async function part2() {
  const data = await parse();

  let edges = new Set<string>();
  for (const [k, v] of data) {
    const connections = new Set<string>([k]);

    for (const value of v) {
      if (connections.has(value)) continue;
      if (!contains(value, connections, data)) continue;

      connections.add(value);
    }

    if (connections.size > edges.size) {
      edges = connections;
    }
  }

  const result = edges.values().toArray().toSorted().join(',');

  console.log('Result:', result);
}

function contains(
  value: string,
  connections: Set<string>,
  data: Map<string, Set<string>>,
) {
  for (const connection of connections) {
    if (!data.get(value)?.has(connection)) {
      return false;
    }
  }

  return true;
}

async function parse() {
  const text = await Deno.readTextFile('assets/2024/task23.txt');
  return text
    .split('\n')
    .map((line) => line.split('-'))
    .reduce((acc, [left, right]) => {
      acc.set(left, (acc.get(left) ?? new Set()).add(right));
      acc.set(right, (acc.get(right) ?? new Set()).add(left));
      return acc;
    }, new Map<string, Set<string>>());
}

if (import.meta.main) {
  main();
}
