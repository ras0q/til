/*!
 * markdown-it-regexp
 * Copyright (c) 2014 Alex Kocharin
 * MIT Licensed
 */

/**
 * @param {RegExp} regexp - Regular expression to match
 * @param {(match: RegExpMatchArray) => string} replacer - Function to replace matched text
 * @returns {(md: MarkdownIt) => void} - Returns a new plugin instance
 */
export default function RegexpPlugin(regexp, replacer) {
  let counter = 0;

  return function (md) {
    const flags = (regexp.global     ? 'g' : '')
                + (regexp.multiline  ? 'm' : '')
                + (regexp.ignoreCase ? 'i' : '')
    const parsedRegexp = RegExp('^' + regexp.source, flags)
    const id = 'regexp-' + counter++

    md.inline.ruler.push(id, (state, silent) => {
      // slowwww... maybe use an advanced regexp engine for this
      const match = parsedRegexp.exec(state.src.slice(state.pos))
      if (!match) return false

      if (state.pending) {
        state.pushPending();
      }

      // valid match found, now we need to advance cursor
      const oldPos = state.pos
      state.pos += match[0].length

      // don't insert any tokens in silent mode
      if (silent) return true

      const token = state.push(id, '', 0)
      token.meta = { match: match }
      token.position = oldPos
      token.size = match[0].length

      return true
    })

    md.renderer.rules[id] = (tokens, idx, _options, _env) => {
      return replacer(tokens[idx].meta.match)
    }
  }
}
