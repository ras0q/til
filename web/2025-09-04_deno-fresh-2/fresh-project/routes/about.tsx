import { Countdown } from "../islands/Countdown.tsx";
import { define } from "../utils.ts";

export default define.page((ctx) => {
  return (
    <main>
      <h1>About</h1>
      <p>This is the about page.</p>
      <p>Shared value: {ctx.state.shared}</p>
      <Countdown />
    </main>
  );
});
