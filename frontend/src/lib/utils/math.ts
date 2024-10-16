export function clamp(input: number, min: number, max: number): number {
  return input < min ? min : input > max ? max : input;
}

export function map(value: number, in_min: number, in_max: number, out_min: number, out_max: number): number {
  const mapped = ((value - in_min) * (out_max - out_min)) / (in_max - in_min) + out_min;

  return clamp(mapped, out_min, out_max);
}

export function rate(numerator: number, denominator: number, decimals: number = 2): number {
  return denominator > 0 ? parseFloat(((numerator / denominator) * 100).toFixed(decimals)) : 0;
}
