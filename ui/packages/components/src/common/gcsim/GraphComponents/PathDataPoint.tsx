import type { MutableRefObject } from "react";

type PathDataPointProps = {
	cx: number;
	x: number;
	fill: string;
	path: MutableRefObject<SVGPathElement | null>;
};

export const PathDataPoint = (props: PathDataPointProps) => {
	if (!props.path.current) {
		return null;
	}

	const y = getPathYFromX(props.x, props.path.current);

	return (
		<g>
			<circle
				cx={props.cx}
				cy={y + 1}
				r={4}
				fill="var(--g-bg)"
				fillOpacity={0.1}
				stroke="var(--g-bg)"
				strokeOpacity={0.1}
				strokeWidth={2}
				pointerEvents="none"
			/>
			<circle
				cx={props.cx}
				cy={y}
				r={4}
				fill={props.fill}
				pointerEvents="none"
				stroke="var(--g-surface)"
				strokeWidth={2}
			/>
		</g>
	);
};

function getPathYFromX(
	x: number,
	path: SVGPathElement,
	error?: number,
): number {
	error = error || 0.01;
	const maxIterations = 10;

	let lengthStart = 0;
	let lengthEnd = path.getTotalLength();
	let point = path.getPointAtLength((lengthEnd + lengthStart) / 2);
	let iterations = 0;

	while (x < point.x - error || x > point.x + error) {
		const midpoint = (lengthStart + lengthEnd) / 2;

		point = path.getPointAtLength(midpoint);

		if (x < point.x) {
			lengthEnd = midpoint;
		} else {
			lengthStart = midpoint;
		}

		iterations += 1;
		if (maxIterations < iterations) {
			break;
		}
	}
	return point.y;
}
