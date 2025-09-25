import { RouteConfig } from "fresh";
import { define } from "../../../utils.ts";
import { Partial } from "fresh/runtime";

// We only want to render the content, so disable
// the `_app.tsx` template as well as any potentially
// inherited layouts
export const config: RouteConfig = {
  skipAppWrapper: true,
  skipInheritedLayouts: true,
};

export default define.page(async (ctx) => {
  const apiResponse = await fetch(
    new URL(`/api/${ctx.params.id}`, ctx.req.url),
  );
  const content = await apiResponse.text();

  // Only render the new content
  return (
    <Partial name="docs-content">
      {content ? <div>{content}</div> : <div>loading...</div>}
    </Partial>
  );
});
