type Display = { color: string; icon: string };

const DISPLAY: Record<string, Display> = {
	action: { color: "#c2410c", icon: "play_arrow" },
	damage: { color: "#2563eb", icon: "local_fire_department" },
	energy: { color: "#059669", icon: "local_cafe" },
	warning: { color: "#dc2626", icon: "circle" },
	status: { color: "#a21caf", icon: "iso" },
	cooldown: { color: "#0d9488", icon: "circle" },
	element: { color: "#4f46e5", icon: "bolt" },
	shield: { color: "#0369a1", icon: "shield" },
	construct: { color: "#78716c", icon: "apartment" },
	player: { color: "#64748b", icon: "circle" },
	user: { color: "#5f7161", icon: "comment" },
	heal: { color: "#16a34a", icon: "healing" },
	hurt: { color: "#991b1b", icon: "coronavirus" },
	pre_damage_mods: { color: "#818cf8", icon: "dynamic_form" },
	icd: { color: "#475569", icon: "timer" },
	calc: { color: "#be185d", icon: "calculate" },
	snapshot: { color: "#6366f1", icon: "photo_camera" },
	character: { color: "#7c8796", icon: "person" },
	weapon: { color: "#9a3412", icon: "circle" },
	enemy: { color: "#7f1d1d", icon: "mood_bad" },
	artifact: { color: "#a16207", icon: "circle" },
	debug: { color: "#52525b", icon: "circle" },
	sim: { color: "#d97706", icon: "circle" },
	hitlag: { color: "#a27b5c", icon: "sports_martial_arts" },
};

const FALLBACK: Display = { color: "#6b7280", icon: "circle" };

export function display(type: string): Display {
	return DISPLAY[type] ?? FALLBACK;
}
