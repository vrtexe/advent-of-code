import { customParseInt } from '../common.ts';

type Comparator<T> = (left: T, right: T) => number;

type Rule = {
  left: number;
  right: number;
};

type Data = {
  rules: Rule[];
  lists: number[][];
};

function main() {
  part2();
}

export async function part1() {
  const { rules, lists } = await parse();
  const groupedRules = Object.groupBy(rules, ({ left }) => left);

  let result = 0;

  for (const list of lists) {
    const applyingRules = extractApplyingRule(list, groupedRules);
    if (validateRules(list, applyingRules)) {
      result += extractMiddleNumber(list);
    }
  }

  console.log('Result:', result);
}

export async function part2() {
  const { rules, lists } = await parse();
  const groupedRules = Object.groupBy(rules, ({ left }) => left);

  let result = 0;

  for (const list of lists) {
    const applyingRules = extractApplyingRule(list, groupedRules);
    if (!validateRules(list, applyingRules)) {
      const comparator = createRuleComparator(applyingRules);
      const reorderedList = reorderByRules(list, comparator);
      result += extractMiddleNumber(reorderedList);
    }
  }

  console.log('Result:', result);
}

function rulesToMap(rules: Rule[]) {
  return new Map(rules.map((rule) => [ruleKeyOf(rule), rule]));
}

function ruleKeyOf(rule: Rule) {
  return `${rule.left}|${rule.right}`;
}

function reorderByRules(list: number[], comparator: Comparator<number>) {
  return list.toSorted(comparator);
}

function createRuleComparator(rules: Rule[]): Comparator<number> {
  const mappedRules = rulesToMap(rules);

  return (left: number, right: number) => {
    if (mappedRules.has(`${left}|${right}`)) return -1;
    if (mappedRules.has(`${right}|${left}`)) return 1;
    return 0;
  };
}

function extractMiddleNumber(list: number[]) {
  return list[(list.length - 1) / 2];
}

function validateRules(list: number[], rules: Rule[]) {
  const orderMap = Object.fromEntries(list.map((v, i) => [v, i]));

  for (const rule of rules) {
    if (orderMap[rule.left] > orderMap[rule.right]) {
      return false;
    }
  }

  return true;
}

function extractApplyingRule(
  list: number[],
  rules: Partial<Record<number, Rule[]>>,
) {
  const items = new Set(list);

  const result: Rule[] = [];

  for (const page of list) {
    if (!rules[page]) continue;
    for (const rule of rules[page]) {
      if (items.has(rule.right)) {
        result.push(rule);
      }
    }
  }

  return result;
}

type Section = (typeof Section)[keyof typeof Section];
const Section = {
  Rule: 'RULE',
  List: 'LIST',
};

async function parse(): Promise<Data> {
  const text = await Deno.readTextFile('src/data/task5.txt');

  let section = Section.Rule;
  const rules: Rule[] = [];
  const lists: number[][] = [];
  for (const line of text.split('\n')) {
    if (!line.trim()) {
      section = Section.List;
      continue;
    }

    if (section === Section.Rule) {
      const [left, right] = line.split('|');
      rules.push({ left: parseInt(left), right: parseInt(right) });
    }

    if (section === Section.List) {
      lists.push(line.split(',').map(customParseInt));
    }
  }

  return { rules, lists };
}

if (import.meta.main) {
  main();
}
