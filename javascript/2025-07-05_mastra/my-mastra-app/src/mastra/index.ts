import { Mastra } from "@mastra/core/mastra";
import { PinoLogger } from "@mastra/loggers";
import { LibSQLStore } from "@mastra/libsql";
import { weatherWorkflow } from "./workflows/weather-workflow.ts";
import { weatherAgent } from "./agents/weather-agent.ts";
import process from "node:process";
import readline from "node:readline/promises";

export const mastra = new Mastra({
  workflows: { weatherWorkflow },
  agents: { weatherAgent },
  storage: new LibSQLStore({
    // stores telemetry, evals, ... into memory storage, if it needs to persist, change to file:../mastra.db
    url: ":memory:",
  }),
  logger: new PinoLogger({
    name: "Mastra",
    level: "info",
  }),
});

const agent = mastra.getAgent("weatherAgent");

const abortController = new AbortController();
const stream = await agent.stream([
  {
    role: "user",
    content:
      "Hello, tell me today's weathers in Tokyo, Kyoto, Osaka, Hokkai-do",
  },
], {
  maxSteps: 10,
  onStepFinish: (stepResult) => {
    console.log("Step completed:", JSON.stringify(stepResult));
  },
  onFinish: (result) => {
    console.log("Stream complete:", JSON.stringify(result));
  },
  abortSignal: abortController.signal,
});

console.log("Agent:");

for await (const part of stream.fullStream) {
  console.log("\n");
  console.log("part:", JSON.stringify(part));
  if (part.type === "tool-call") {
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });
    const answer = await rl.question(
      `Can I execute a tool (${part.toolName})? [y/N]`,
    );
    rl.close();
    if (answer === "N") {
      abortController.abort();
      break;
    }
  } else if (part.type === "text-delta") {
    process.stdout.write(part.textDelta);
  }
}
