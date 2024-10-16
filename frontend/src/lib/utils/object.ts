export function clone<T>(object: T): T {
  if (object === undefined || object === null) return object;
  return JSON.parse(JSON.stringify(object));
}

export function compare<A, B>(a: A, b: B): boolean {
  return JSON.stringify(a) === JSON.stringify(b);
}

export function major(object: Record<string, number>): string {
  let majorKey = "";
  let majorValue = Number.MIN_SAFE_INTEGER;

  Object.entries(object).forEach(([key, value]) => {
    if (value > majorValue || (value === majorValue && key > majorKey)) {
      majorKey = key;
      majorValue = value;
    }
  });

  return majorKey;
}
