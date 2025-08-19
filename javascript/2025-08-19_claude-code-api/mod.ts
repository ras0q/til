import { Hono } from "hono";
import type { paths } from "./openai.d.ts";
import { createHonoOpenApiRouter } from "openapi-ts-router";
import { vValidator } from "validation-adapters/valibot";
import * as v from "valibot";

const app = new Hono();
export const openApiRouter = createHonoOpenApiRouter<paths>(app);

openApiRouter.get("/chat/completions", {
  queryValidator: vValidator(
    v.object({}) satisfies v.GenericSchema<
      NonNullable<paths["/chat/completions"]["get"]["parameters"]["query"]>
    >,
  ),
  handler: (c) => {
    const query = c.req.valid("query");
    console.log("query", query);

    return c.json(
      {
        "object": "list",
        "data": [
          {
            "id": "string",
            "choices": [
              {
                "finish_reason": "stop",
                "index": 0,
                "message": {
                  "content": "string",
                  "refusal": "string",
                  "tool_calls": [
                    {
                      "id": "string",
                      "type": "function",
                      "function": {
                        "name": "string",
                        "arguments": "string",
                      },
                    },
                    {
                      "id": "string",
                      "type": "custom",
                      "custom": {
                        "name": "string",
                        "input": "string",
                      },
                    },
                  ],
                  "annotations": [
                    {
                      "type": "url_citation",
                      "url_citation": {
                        "end_index": 0,
                        "start_index": 0,
                        "url": "string",
                        "title": "string",
                      },
                    },
                  ],
                  "role": "assistant",
                  "audio": {
                    "id": "string",
                    "expires_at": 0,
                    "data": "string",
                    "transcript": "string",
                  },
                },
                "logprobs": {
                  "content": [
                    {
                      "token": "string",
                      "logprob": 0,
                      "bytes": [
                        0,
                      ],
                      "top_logprobs": [
                        {
                          "token": "string",
                          "logprob": 0,
                          "bytes": [
                            0,
                          ],
                        },
                      ],
                    },
                  ],
                  "refusal": [
                    {
                      "token": "string",
                      "logprob": 0,
                      "bytes": [
                        0,
                      ],
                      "top_logprobs": [
                        {
                          "token": "string",
                          "logprob": 0,
                          "bytes": [
                            0,
                          ],
                        },
                      ],
                    },
                  ],
                },
              },
            ],
            "created": 0,
            "model": "string",
            "service_tier": "auto",
            "object": "chat.completion",
            "usage": {
              "completion_tokens": 0,
              "prompt_tokens": 0,
              "total_tokens": 0,
              "completion_tokens_details": {
                "accepted_prediction_tokens": 0,
                "audio_tokens": 0,
                "reasoning_tokens": 0,
                "rejected_prediction_tokens": 0,
              },
              "prompt_tokens_details": {
                "audio_tokens": 0,
                "cached_tokens": 0,
              },
            },
          },
        ],
        "first_id": "string",
        "last_id": "string",
        "has_more": true,
      } as const,
    );
  },
});
