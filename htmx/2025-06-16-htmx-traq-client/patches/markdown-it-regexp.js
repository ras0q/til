/*!
 * markdown-it-regexp
 * Copyright (c) 2014 Alex Kocharin
 * MIT Licensed
 */

// `MarkdownIt` import is for type hinting in TypeScript, not needed in JavaScript directly.

/**
 * Counter for multiple usages of the plugin to ensure unique IDs.
 */
let counter = 0;

class MarkdownItRegexp {
  /**
   * Creates an instance of MarkdownItRegexp.
   */
  constructor(regexp, replacer) {
    // Clone regexp with all the flags
    const flags =
      (regexp.global ? "g" : "") +
      (regexp.multiline ? "m" : "") +
      (regexp.ignoreCase ? "i" : "");

    // Ensure the regex matches from the start of the current position in the state.src
    this.regexp = new RegExp("^" + regexp.source, flags);

    // Copy init options
    this.replacer = replacer;

    // This plugin can be inserted multiple times,
    // so we're generating a unique name for it
    this.id = "regexp-" + counter;
    counter++;
  }

  /**
   * Function that registers the plugin with markdown-it.
   */
  init(md) {
    md.inline?.ruler.push(this.id, this.parse.bind(this));
    if (!md.renderer) {
      md.renderer = { rules: {} };
    }
    md.renderer.rules[this.id] = this.render.bind(this);
  }

  /**
   * Parser function for the markdown-it inline rule.
   * @param {any} state The current markdown-it inline state.
   * @param {boolean} silent If true, the parser should not insert tokens.
   * @returns {boolean} True if a match was found and processed, false otherwise.
   */
  parse(state, silent) {
    // slowwww... maybe use an advanced regexp engine for this
    const match = this.regexp.exec(state.src.slice(state.pos));
    if (!match) {
      return false;
    }

    // Original code had `if (state.pending) { state.pushPending(); }`.
    // As noted previously, `state.pushPending()` is not a standard markdown-it API call in modern versions.
    // This line is kept for direct translation but may cause issues with standard markdown-it if `state.pushPending` is undefined.
    if (state.pending) {
      state.pushPending();
    }

    // Valid match found, now we need to advance cursor
    const oldPos = state.pos;
    state.pos += match[0].length;

    // Don't insert any tokens in silent mode
    if (silent) return true;

    const token = state.push(this.id, "", 0);
    token.meta = { match: match };
    token.position = oldPos; // Non-standard markdown-it token property, but kept for direct translation
    token.size = match[0].length; // Non-standard markdown-it token property, but kept for direct translation

    return true;
  }

  /**
   * Renderer function for the custom token.
   */
  render(tokens, idx, options, env) {
    const token = tokens[idx];
    if (token && token.meta && token.meta.match) {
      return this.replacer(token.meta.match);
    }
    return ""; // Should ideally not happen if parsing is correct
  }
}

// Export the class directly for use as a markdown-it plugin
export default function (md, regexp, replacer) {
  const plugin = new MarkdownItRegexp(regexp, replacer);
  return (md, options) => {
    plugin.init(md);
  }
}
