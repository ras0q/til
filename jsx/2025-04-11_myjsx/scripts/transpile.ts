import { transpile } from "@deno/emit";

const url = new URL(Deno.args[0], import.meta.url);
console.log(url);
const result = await transpile(url, {});

const code = result.get(url.href);
console.log(code);
