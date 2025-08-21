export function customParseInt(value: string) {
  return parseInt(value);
}

export function sum(prev: number, next: number) {
  return prev + next;
}

export function sortAscending(left: number, right: number) {
  return left - right;
}

export function multiply(prev: number, next: number) {
  return prev * next;
}

export function minFunc<T>(compare: (left: T, right: T) => number) {
  return (result: T | undefined, current: T) =>
    !result || compare(current, result) < 0 ? current : result;
}

export function maxFunc<T>(compare: (left: T, right: T) => number) {
  return (result: T, current: T) =>
    !result || compare(current, result) > 0 ? current : result;
}

export function min(result: number, current: number) {
  return !result || current - result < 0 ? current : result;
}

export function max(result: number, current: number) {
  return !result || current - result > 0 ? current : result;
}
