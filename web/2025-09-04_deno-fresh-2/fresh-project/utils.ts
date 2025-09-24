import { createDefine } from "fresh";

// This specifies the type of "ctx.state" which is used to share
// data among middlewares, layouts and routes.
export interface State {
  shared: string;
}

export const define = createDefine<State>();

export async function loadContent(id: string) {
  await new Promise((resolve) => setTimeout(() => resolve(0), 100));

  return `Hello, ${id}!`;
}
