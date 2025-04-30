import { globToRegExp } from "jsr:@std/path@1.0.9";
import { walk } from "jsr:@std/fs@1.0.17";

const glob = [
  globToRegExp(
    "*/[[:digit:]][[:digit:]][[:digit:]][[:digit:]]-[[:digit:]][[:digit:]]-[[:digit:]][[:digit:]]_*/",
  ),
];

const tilEntries = (await Array.fromAsync(walk(".", { match: glob })))
  .filter((e) => e.isDirectory)
  .toSorted((a, b) => a.name === b.name ? 0 : a.name > b.name ? 1 : -1);

const content = `# Today I Learned (TIL)

## Index (${tilEntries.length} entries, Newest first)

${
  tilEntries
    .map((e) => `- [${e.path}](./${e.path})`)
    .join("\n")
}`;

console.log(content);
