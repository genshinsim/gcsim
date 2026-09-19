import { ChevronDownIcon, ChevronUpIcon } from "lucide-react";
import * as React from "react";

import { cn } from "../../lib/utils";

import { Button } from "./button";
import { Input } from "./input";

type NumberInputProps = Omit<
	React.ComponentProps<"input">,
	"value" | "onChange" | "type" | "min" | "max" | "step"
> & {
	value: number;
	onValueChange: (value: number) => void;
	min?: number;
	max?: number;
	step?: number;
};

function clamp(value: number, min?: number, max?: number) {
	if (min != null && value < min) return min;
	if (max != null && value > max) return max;
	return value;
}

function toFiniteNumber(text: string) {
	const trimmed = text.trim();
	if (trimmed === "") return null;
	const value = Number(trimmed);
	return Number.isFinite(value) ? value : null;
}

function NumberInput({
	value,
	onValueChange,
	min,
	max,
	step = 1,
	className,
	disabled,
	onBlur,
	onKeyDown,
	...props
}: NumberInputProps) {
	const [text, setText] = React.useState(() => String(value));

	React.useEffect(() => {
		setText((prev) => {
			const parsed = toFiniteNumber(prev);
			return parsed !== null && clamp(parsed, min, max) === value
				? prev
				: String(value);
		});
	}, [value, min, max]);

	const commit = (next: number) => {
		const clamped = clamp(next, min, max);
		setText(String(clamped));
		onValueChange(clamped);
	};

	const stepBy = (direction: 1 | -1) =>
		commit((Number.isFinite(value) ? value : 0) + direction * step);

	return (
		<div
			data-slot="number-input"
			className={cn("flex w-full items-stretch", className)}
		>
			<Input
				type="text"
				inputMode="decimal"
				disabled={disabled}
				value={text}
				className="rounded-r-none"
				onChange={(event) => {
					const raw = event.target.value;
					setText(raw);
					const parsed = toFiniteNumber(raw);
					if (parsed !== null) onValueChange(clamp(parsed, min, max));
				}}
				onBlur={(event) => {
					commit(toFiniteNumber(text) ?? value);
					onBlur?.(event);
				}}
				onKeyDown={(event) => {
					if (event.key === "ArrowUp") {
						event.preventDefault();
						stepBy(1);
					} else if (event.key === "ArrowDown") {
						event.preventDefault();
						stepBy(-1);
					}
					onKeyDown?.(event);
				}}
				{...props}
			/>
			<div className="flex flex-col">
				<StepperButton
					direction="up"
					disabled={disabled || (max != null && value >= max)}
					onClick={() => stepBy(1)}
				/>
				<StepperButton
					direction="down"
					disabled={disabled || (min != null && value <= min)}
					onClick={() => stepBy(-1)}
				/>
			</div>
		</div>
	);
}

function StepperButton({
	direction,
	className,
	...props
}: React.ComponentProps<typeof Button> & { direction: "up" | "down" }) {
	return (
		<Button
			type="button"
			variant="ghost"
			tabIndex={-1}
			data-slot="number-input-stepper"
			className={cn(
				"h-auto min-w-0 flex-1 rounded-none border border-l-0 border-g-line px-2 py-0 text-g-ink-mute focus-visible:relative focus-visible:z-10",
				direction === "up"
					? "rounded-tr-g-input border-b-0"
					: "rounded-br-g-input",
				className,
			)}
			{...props}
		>
			{direction === "up" ? (
				<ChevronUpIcon className="size-3" />
			) : (
				<ChevronDownIcon className="size-3" />
			)}
		</Button>
	);
}

export type { NumberInputProps };
export { NumberInput };
