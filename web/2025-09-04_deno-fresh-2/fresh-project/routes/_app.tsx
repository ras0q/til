import { define } from "../utils.ts";
import { Partial } from "fresh/runtime";

export default define.page(function App({ Component }) {
  return (
    <html>
      <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>fresh-project</title>
      </head>
      <body f-client-nav>
        <aside class="flex gap-8">
          <a href="/">Home</a>
          <a href="/about">About</a>
          <a href="/docs/routes">Routes</a>
        </aside>
        <Partial name="body">
          <Component />
        </Partial>
      </body>
    </html>
  );
});
