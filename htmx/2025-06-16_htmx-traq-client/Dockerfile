FROM denoland/deno:2.3.6

RUN apt-get update && apt-get install -y patch && rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /app

# Copy source
COPY deno.jsonc deno.lock main.tsx schema.d.ts ./
COPY ./patches ./patches

# Compile the main app
RUN deno cache main.tsx && deno task patch

EXPOSE 8000

# Run the app
ENTRYPOINT ["/app/main.tsx"]
