export function trimPrefix(str: string, prefix: string): string {
  if (str.startsWith(prefix)) {
    return str.slice(prefix.length, str.length);
  }

  return str;
}

export function trimSuffix(str: string, suffix: string): string {
  if (str.endsWith(suffix)) {
    return str.slice(0, str.length - suffix.length);
  }

  return str;
}

export function trim(str: string, fix: string): string {
  return trimPrefix(trimSuffix(str, fix), fix);
}

export function capitalize(str: string): string {
  return str[0].toUpperCase() + str.slice(1).toLowerCase();
}

export function titlelize(str: string): string {
  return str
    .split(" ")
    .map((word) => capitalize(word))
    .join(" ");
}

export function unitize(value: number): string {
  if (value < 1000) {
    return `${value}`;
  }

  let div = 1000;
  let exp = 0;
  for (let n = value / 1000; n >= 1000; n /= 1000) {
    div *= 1000;
    exp++;
  }

  const number = trimSuffix((value / div).toFixed(1), ".0");
  const exponent = "KMGTPE"[exp];

  return `${number}${exponent}`;
}

const PUNCTUATION_SET: string[] = [
  "!",
  '"',
  "#",
  "$",
  "%",
  "&",
  "'",
  "(",
  ")",
  "*",
  "+",
  ",",
  "-",
  ".",
  "/",
  ":",
  ";",
  "?",
  "@",
  "[",
  "\\",
  "]",
  "^",
  "_",
  "`",
  "{",
  "|",
  "}",
  "~",
];

export function initials(str: string): string {
  const parts = trim(str.replaceAll(".", ""), " ").split(" ");

  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  if (PUNCTUATION_SET.includes(parts[1][0])) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  return (parts[0][0] + parts[1][0]).toUpperCase();
}

const RANDOM_LETTER_SET: string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

export function random(length: number): string {
  let result = "";

  for (let i = 0; i < length; i++)
    result += RANDOM_LETTER_SET.charAt(Math.floor(Math.random() * RANDOM_LETTER_SET.length));

  return result;
}

const DETERMINERS_REGEX = /\s?(a few|an|a)\s/g;
const SECONDS_REGEX = /\s?(seconds|second)/g;
const MINUTES_REGEX = /\s?(minutes|minute)/g;
const HOURS_REGEX = /\s?(hours|hour)/g;
const DAYS_REGEX = /\s?(days|day)/g;
const MONTHS_REGEX = /\s?(months|month)/g;
const YEARS_REGEX = /\s?(years|year)/g;

export function simplify(str: string): string {
  str = str.replaceAll(DETERMINERS_REGEX, "1");
  str = str.replaceAll(SECONDS_REGEX, "s");
  str = str.replaceAll(MINUTES_REGEX, "min");
  str = str.replaceAll(HOURS_REGEX, "hr");
  str = str.replaceAll(DAYS_REGEX, "d");
  str = str.replaceAll(MONTHS_REGEX, "mo");
  str = str.replaceAll(YEARS_REGEX, "y");

  return str;
}
