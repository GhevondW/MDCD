// Keep emoji as native unicode (no twemoji CDN images) so the exported
// HTML has zero network dependencies and presents fully offline.
//
// html: true is required — the deck uses raw <div style="..."> for the
// "Three families" column layout and inline <svg class="uml"> diagrams.
// Without it, marp-cli silently strips style attributes (tags survive,
// styling doesn't) instead of failing loudly, so don't drop this.
module.exports = {
  html: true,
  options: {
    emoji: { shortcode: false, unicode: false },
  },
}
