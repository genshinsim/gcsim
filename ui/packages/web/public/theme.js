// Runs before the bundle so a saved theme paints first. The key is THEME_KEY
// in ui/src/stores/prefs.ts; an unknown id matches no palette and shows Cryo.
try {
	const theme = localStorage.getItem("gcsim-theme");
	if (theme) document.documentElement.dataset.theme = theme;
} catch {}
