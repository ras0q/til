import { Partial } from "fresh/runtime";
import { define } from "../../utils.ts";

export default define.page(() => {
  return (
    <div>
      <aside>
        <a href="/docs/page1" f-partial="/partials/docs/page1">Page 1</a>
        <a href="/docs/page2" f-partial="/partials/docs/page2">Page 2</a>
      </aside>
      <Partial name="docs-content">Page is not selected.</Partial>
    </div>
  );
});
