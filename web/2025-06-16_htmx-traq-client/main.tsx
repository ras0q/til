#!/usr/bin/env -S deno serve --allow-net --allow-env

import { Hono } from "hono";
import type { FC } from "hono/jsx";
import { HTTPException } from "hono/http-exception";
import { authHandler, initAuthConfig, verifyAuth } from "@hono/auth-js";
import { streamSSE } from "hono/streaming";
import type { components, paths } from "./schema.d.ts";
import createClient from "openapi-fetch";
import { traQMarkdownIt } from "@traptitech/traq-markdown-it";

const apiBaseURL = "https://q.trap.jp/api/v3";

const globalStyle = `
.message {
  display: flex;
  gap: 1rem;
  align-items: top;
  margin-bottom: 1rem;
  padding: 0.5rem;
}

.message-header {
  display: flex;
  gap: 0.5rem;

  > span {
    color: #666;
  }
}

.message-content {
  h1,h2,h3,p { margin: 0.5rem 0; }
}

.user-icon {
  border-radius: 50%;
  margin-top: 0.25rem;
}
`;

const Layout: FC = (props) => {
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>traQ Viewer</title>
        <script
          src="https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js"
          integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"
          crossorigin="anonymous"
        />
        <script
          src="https://unpkg.com/htmx-ext-sse@2.2.3/dist/sse.min.js"
          integrity="sha384-Y4gc0CK6Kg+hmulDc6rZPJu0tqvk7EWlih0Oh+2OkAi1ZDlCbBDCQEE2uVk472Ky"
          crossorigin="anonymous"
        />
        {/* <script src="https://cdn.jsdelivr.net/npm/@unocss/runtime" /> */}
        <style>{globalStyle}</style>
      </head>
      <body>{props.children}</body>
    </html>
  );
};

const MessageComponent: FC<
  {
    message: components["schemas"]["Message"];
    userMap: Map<string, components["schemas"]["User"]>;
    md: traQMarkdownIt;
  }
> = (
  { message, userMap, md },
) => {
  const user = userMap.get(message.userId);
  const username = user?.name || "unknown";
  const userDisplayName = user?.displayName || username;
  const createdAt = new Date(message.createdAt).toLocaleString("ja-JP", {
    timeZone: "Asia/Tokyo",
  });
  const markdownContent: string = md.render(message.content || "").renderedText;

  return (
    <div class="message">
      <img
        class="user-icon"
        src={`${apiBaseURL}/public/icon/${username}`}
        loading="lazy"
        width="40px"
        height="40px"
      />
      <div>
        <div class="message-header">
          <strong>{userDisplayName}</strong>
          <span>@{username}</span>
          <span>{createdAt}</span>
        </div>
        <div
          class="message-content"
          dangerouslySetInnerHTML={{ __html: markdownContent }}
        />
        {message.stamps.length > 0 && (
          <div>
            {message.stamps.map((stamp) => (
              <img
                key={stamp.stampId}
                src={`${apiBaseURL}/public/icons/traP`}
                alt="TODO: fill stamp name"
                title="traQ"
                width="32"
                height="32"
                loading="lazy"
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

const app = new Hono();

app.use(
  "*",
  initAuthConfig((c) => ({
    basePath: "/auth",
    secret: c.env.AUTH_SECRET,
    providers: [{
      id: "traq",
      name: "traQ",
      type: "oauth",
      authorization: {
        url: `${apiBaseURL}/oauth2/authorize`,
        params: {
          scope: "read write",
        },
      },
      token: `${apiBaseURL}/oauth2/token`,
      userinfo: `${apiBaseURL}/users/me`,
      profile(profile) {
        return {
          id: profile.id,
          name: profile.name,
          image: `${apiBaseURL}/public/icon/${profile.name}`,
        };
      },
    }],
    callbacks: {
      session({ session, token }) {
        return {
          ...session,
          user: { ...session.user, accessToken: token.accessToken },
        };
      },
      jwt({ token, account }) {
        if (account) {
          return {
            ...token,
            accessToken: account?.access_token,
          };
        }
        return token;
      },
    },
  })),
);
app.use("/auth/*", authHandler());
app.use("*", verifyAuth());

app.get("/", (c) => {
  const user = c.get("authUser");
  if (!user) {
    return c.redirect("/auth/signin");
  }

  return c.html(
    <Layout>
      <h1>traQ Viewer</h1>
      <pre>{JSON.stringify(user, null, 2)}</pre>
      <a href="/channels/general">/channels/general</a>
    </Layout>,
  );
});

app.get("/channels/:path{.+}", (c) => {
  const channelPath = c.req.param("path");

  return c.html(
    <Layout>
      <h1>#{channelPath}</h1>
      <div
        hx-ext="sse"
        sse-connect={`/api/channels/${channelPath}/stream`}
      >
        <div sse-swap="message" hx-swap="beforeend" />
        <div sse-swap="error" hx-swap="beforeend" />
      </div>
    </Layout>,
  );
});

app.get("/api/channels/:path{.+}/stream", async (c) => {
  const user = c.get("authUser");
  const accessToken = user?.token?.accessToken as string | undefined;
  if (!accessToken) {
    return c.json("Unauthorized", 401);
  }

  const apiClient = createClient<paths>({
    baseUrl: apiBaseURL,
    headers: {
      Authorization: `Bearer ${accessToken}`,
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });

  const channelPath = c.req.param("path");
  const { data: channelsByPath, error: getChannelByPathError } = await apiClient
    .GET(
      "/channels",
      {
        params: { query: { path: channelPath } },
      },
    );
  if (getChannelByPathError) return c.json("failed to fetch channel", 500);
  if (channelsByPath.public.length === 0) {
    return c.json("Channel not found", 404);
  }
  const channel = channelsByPath.public[0];
  if (!channel) return c.json("Channel not found", 404);

  const { data: channels, error: getChannelsError } = await apiClient.GET(
    "/channels",
  );
  if (getChannelsError) return c.json("failed to fetch channels", 500);
  if (!channels) return c.json("Channels not found", 404);
  const channelMap = new Map(channels.public.map((c) => [c.id, c]));

  const { data: users, error: getUsersError } = await apiClient.GET("/users", {
    params: {
      query: {
        "include-suspended": true,
      },
    },
  });
  if (getUsersError) return c.json("failed to fetch users", 500);
  if (!users) return c.json("Users not found", 404);
  const userMap = new Map(users.map((u) => [u.id, u]));

  const { data: userGroups, error: getUserGroupsError } = await apiClient.GET(
    "/groups",
  );
  if (getUserGroupsError) return c.json("failed to fetch user groups", 500);
  if (!userGroups) return c.json("User groups not found", 404);
  const userGroupMap = new Map(userGroups.map((g) => [g.id, g]));

  const { data: stamps, error: getStampsError } = await apiClient.GET(
    "/stamps",
  );
  if (getStampsError) return c.json("failed to fetch stamps", 500);
  if (!stamps) return c.json("Stamps not found", 404);
  const stampMap = new Map(stamps.map((s) => [s.name, s]));

  const md = new traQMarkdownIt(
    {
      getChannel: () => channel,
      getMe: () => undefined,
      getUser: (userId) => userMap.get(userId),
      getUserGroup: (userGroupId) => userGroupMap.get(userGroupId),
      getStampByName: (name) => {
        const stamp = stampMap.get(name);
        console.log(stamp);
        return stamp;
      },
      getUserByName: (username) => users.find((u) => u.name === username),
      generateUserHref: (userId) => `/users/${userId}`,
      generateUserGroupHref: (userGroupId) => `/groups/${userGroupId}`,
      generateChannelHref: (channelId) => {
        let channel = channelMap.get(channelId);
        let path = "";
        while (channel?.parentId) {
          path = `/${channel.name}${path}`;
          channel = channelMap.get(channel.parentId);
        }
        return `/channels/${channel?.name}${path}`;
      },
    },
    [],
    "https://q.trap.jp",
  );

  let lastMessage: components["schemas"]["Message"] | null = null;
  return streamSSE(c, async (stream) => {
    while (true) {
      const order = "asc";
      const limit = 200;
      const since = lastMessage
        ? lastMessage.createdAt
        : new Date(new Date().getTime() - 3 * 60 * 60 * 1000) // 3 hours ago
          .toISOString();

      const { data: messages, error: getMessagesError } = await apiClient.GET(
        "/channels/{channelId}/messages",
        {
          params: {
            path: { channelId: channel.id },
            query: { order, limit, since },
          },
        },
      );
      if (getMessagesError) {
        stream.writeSSE({
          data: JSON.stringify({ error: "Failed to fetch messages" }),
          event: "error",
        });
        break;
      }

      if (!messages || (messages.length === 0 && lastMessage)) {
        continue;
      }

      lastMessage = messages[messages.length - 1];

      const data = (
        <div>
          {messages.map((message) => (
            <MessageComponent
              key={message.id}
              message={message}
              md={md}
              userMap={userMap}
            />
          ))}
        </div>
      );

      stream.writeSSE({
        data: data.toString(),
        event: "message",
      });

      await stream.sleep(10 * 1000);
    }
  });
});

app.onError((err, c) => {
  if (err instanceof HTTPException && err.status === 401) {
    return c.redirect("/auth/signin");
  }

  return c.json(
    {
      error: "Internal Server Error",
      message: err.message,
      stack: err.stack,
    },
    500,
  );
});

export default app;
