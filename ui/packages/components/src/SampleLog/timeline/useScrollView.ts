import { useEffect, useState } from "react";

type View = { left: number; width: number };

export function useScrollView() {
	const [el, setEl] = useState<HTMLDivElement | null>(null);
	const [view, setView] = useState<View>({ left: 0, width: 1000 });
	useEffect(() => {
		if (el == null) {
			return;
		}
		const on = () =>
			setView((v) =>
				v.left === el.scrollLeft && v.width === el.clientWidth
					? v
					: { left: el.scrollLeft, width: el.clientWidth },
			);
		on();
		el.addEventListener("scroll", on, { passive: true });
		const ro = new ResizeObserver(on);
		ro.observe(el);
		return () => {
			el.removeEventListener("scroll", on);
			ro.disconnect();
		};
	}, [el]);
	return { el, ref: setEl, view };
}
