import { RouteConfig } from "fresh";
import { define, loadContent } from "../../../utils.ts";
import { Partial } from "fresh/runtime";

// We only want to render the content, so disable
// the `_app.tsx` template as well as any potentially
// inherited layouts
export const config: RouteConfig = {
  skipAppWrapper: true,
  skipInheritedLayouts: true,
};

export default define.page(async (ctx) => {
  const content = await loadContent(ctx.params.id);

  // Only render the new content
  return (
    <Partial name="docs-content">
      {content ? <div>{content}</div> : <div>loading...</div>}
    </Partial>
  );
});
