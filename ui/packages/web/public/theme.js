try {
	const theme = localStorage.getItem("gcsim-theme");
	if (theme) document.documentElement.dataset.theme = theme;
} catch {}
