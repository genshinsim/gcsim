import { type JSX, useEffect, useRef, useState } from "react";

type Props = {
	children: (width: number) => JSX.Element;
};

// Unlike @visx's ParentSize, children render in normal flow, so a chart that
// sets its own height sizes the container instead of being clipped by
// ParentSize's absolute measurement box.
export const ParentWidth = ({ children }: Props) => {
	const ref = useRef<HTMLDivElement>(null);
	const [width, setWidth] = useState(0);

	useEffect(() => {
		const el = ref.current;
		if (el == null) {
			return;
		}
		const observer = new ResizeObserver((entries) => {
			setWidth(entries[0].contentRect.width);
		});
		observer.observe(el);
		return () => observer.disconnect();
	}, []);

	return (
		<div ref={ref} className="w-full">
			{width > 0 ? children(width) : null}
		</div>
	);
};
