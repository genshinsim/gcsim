import { useLayoutEffect, useRef } from "react";

const MIN_HEIGHT = 280;

function scrollParent(el: HTMLElement): HTMLElement {
	for (let p = el.parentElement; p != null; p = p.parentElement) {
		if (/auto|scroll/.test(getComputedStyle(p).overflowY)) {
			return p;
		}
	}
	return document.documentElement;
}

export function useFitViewport() {
	const ref = useRef<HTMLDivElement>(null);
	useLayoutEffect(() => {
		const el = ref.current;
		if (el == null) {
			return;
		}
		const page = scrollParent(el);
		const resize = () => {
			const pageTop =
				page === document.documentElement
					? 0
					: page.getBoundingClientRect().top;
			const top = el.getBoundingClientRect().top - pageTop + page.scrollTop;
			// overflow the page so scrollHeight measures the content, not the page
			el.style.height = `${page.clientHeight + 1}px`;
			const below = page.scrollHeight - top - el.offsetHeight;
			const fitsUnderAbove = page.clientHeight - top - below;
			const h = Math.max(
				MIN_HEIGHT,
				fitsUnderAbove >= MIN_HEIGHT
					? fitsUnderAbove
					: page.clientHeight - below,
			);
			el.style.height = `${h}px`;
		};
		resize();
		window.addEventListener("resize", resize);
		const ro = new ResizeObserver(resize);
		ro.observe(page);
		ro.observe(document.body);
		return () => {
			window.removeEventListener("resize", resize);
			ro.disconnect();
		};
	}, []);
	return ref;
}
