import anniversary from "./hero/anniversary.jpg";
import bilibili10m from "./hero/bilibili-10m.png";
import christmas2021 from "./hero/christmas-2021.png";
import ganyu from "./hero/ganyu.jpg";
import ganyuLiyueHarbor from "./hero/ganyu-liyue-harbor.png";
import version24 from "./hero/version-2-4.jpg";

export type HeroImage = {
	id: string;
	name: string;
	src: string;
	knobs?: Record<string, string>;
};

export const HERO_IMAGES: HeroImage[] = [
	{ id: "anniversary", name: "Anniversary", src: anniversary },
	{ id: "ganyu", name: "Ganyu", src: ganyu },
	{
		id: "ganyu-liyue-harbor",
		name: "Ganyu · Liyue Harbor",
		src: ganyuLiyueHarbor,
	},
	{ id: "version-2-4", name: "Version 2.4", src: version24 },
	{ id: "bilibili-10m", name: "Bilibili 10M", src: bilibili10m },
	{ id: "christmas-2021", name: "Christmas 2021", src: christmas2021 },
];

export const DEFAULT_HERO_ID = "anniversary";

export function getHero(id: string = DEFAULT_HERO_ID): HeroImage {
	return (
		HERO_IMAGES.find((h) => h.id === id) ??
		HERO_IMAGES.find((h) => h.id === DEFAULT_HERO_ID) ??
		HERO_IMAGES[0]
	);
}
