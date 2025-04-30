<script setup lang="ts">
import type { WindowNostr } from 'nostr-tools/nip07'
import { ref } from 'vue';
import { Relay, type RelayRecord } from 'nostr-tools/relay';
import type { Event } from 'nostr-tools';

declare global {
  interface Window {
    nostr: WindowNostr;
  }
}

const publicKey = ref<string>()
const relays = ref<RelayRecord>()
const events = ref<Event[]>([])
const isStreaming = ref<Record<string, boolean>>({})

const handleLoginButtonClick = async () => {
  publicKey.value = await window.nostr.getPublicKey()
  relays.value = await window.nostr.getRelays()

  for (const [relayURL, { read, write }] of Object.entries(relays.value)) {
    if (!read && !write) continue

    const relay = await Relay.connect(relayURL)
    console.info(`Connected: ${relay.url}`)

    isStreaming.value[relay.url] = false
    relay.subscribe([{
      // https://github.com/nostr-protocol/nips#event-kinds
      kinds: [
        1 // Short Text Note
      ],
      limit: 20,
    }], {
      onevent: (e) => {
        if (isStreaming.value[relay.url]) {
          console.info(`New Event(streaming): ${e.id}`)
          events.value.unshift(e)
        } else {
          console.info(`New Event(stored): ${e.id}`)
          events.value.push(e)
        }
      },
      oneose: () => {
        console.info("Eose")
        isStreaming.value[relay.url] = true
      },
      onclose: (reason) => {
        console.info(`Closing: ${reason}`)
      }
    })
  }
}
</script>

<template>
  <div>
    <button @click="handleLoginButtonClick">Login with NIP-07</button>
    <div>PublicKey: {{ publicKey }}</div>
    <div>Relays: {{ relays }}</div>
    <div>Events:</div>
    <ul>
      <li v-for="event in events" :key="event.id">{{
        `${event.pubkey.substring(0, 6)} (${new
          Date(event.created_at * 1000).toISOString().substring(11, 16)}):
        ${event.content}` }}</li>
    </ul>
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
