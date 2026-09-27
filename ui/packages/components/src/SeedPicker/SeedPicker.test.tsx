import type { model } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (k: string, o?: { p: number }) => (o ? `${k}:${o.p}` : k),
	}),
}));

// Radix Select takes seconds to open under jsdom; a native select exercises
// the same value/onValueChange contract.
vi.mock("@gcsim/primitives", async (importOriginal) => ({
	...(await importOriginal<object>()),
	Select: ({
		value,
		onValueChange,
		children,
	}: {
		value: string;
		onValueChange: (v: string) => void;
		children: React.ReactNode;
	}) => (
		<select value={value} onChange={(e) => onValueChange(e.target.value)}>
			{children}
		</select>
	),
	SelectTrigger: () => null,
	SelectValue: () => null,
	SelectContent: ({ children }: { children: React.ReactNode }) => (
		<>{children}</>
	),
	SelectItem: ({
		value,
		children,
	}: {
		value: string;
		children: React.ReactNode;
	}) => <option value={value}>{children}</option>,
}));

import { namedSeeds, SeedPicker } from "./SeedPicker";

const seeds = {
	sample: "18446744073709551615",
	min: "11",
	max: "22",
	p25: "33",
	p50: "44",
	p75: "55",
};

function renderPicker(
	props: Partial<React.ComponentProps<typeof SeedPicker>> = {},
) {
	const onPick = vi.fn();
	render(<SeedPicker seeds={seeds} value={null} onPick={onPick} {...props} />);
	return onPick;
}

const choose = (label: string) =>
	userEvent.selectOptions(screen.getByRole("combobox"), label);

const generate = () =>
	userEvent.click(screen.getByRole("button", { name: /viewer.generate/ }));

describe("SeedPicker", () => {
	it("picks the sample seed by default, keeping the full uint64 string", async () => {
		const onPick = renderPicker();
		await generate();
		expect(onPick).toHaveBeenCalledWith("18446744073709551615");
	});

	it("picks the chosen named seed", async () => {
		const onPick = renderPicker();
		await choose("viewer.seed_p:75");
		await generate();
		expect(onPick).toHaveBeenCalledWith("55");
	});

	it("starts on the option matching the current value", async () => {
		const onPick = renderPicker({ value: "22" });
		expect(screen.getByRole("combobox")).toHaveValue("max");
		await generate();
		expect(onPick).toHaveBeenCalledWith("22");
	});

	it("accepts a custom seed", async () => {
		const onPick = renderPicker();
		await choose("viewer.seed_custom");
		const input = screen.getByRole("textbox");
		await userEvent.clear(input);
		await userEvent.type(input, "9007199254740993");
		await generate();
		expect(onPick).toHaveBeenCalledWith("9007199254740993");
	});

	it("starts on custom with the value filled in when it matches no named seed", () => {
		renderPicker({ value: "12345" });
		expect(screen.getByRole("combobox")).toHaveValue("custom");
		expect(screen.getByRole("textbox")).toHaveValue("12345");
	});

	it("disables statistic seeds while the run is still going", async () => {
		renderPicker({ running: true });
		expect(
			screen.getByRole("button", { name: /viewer.generate/ }),
		).toBeEnabled();
		await choose("viewer.seed_min");
		expect(
			screen.getByRole("button", { name: /viewer.generate/ }),
		).toBeDisabled();
	});
});

describe("namedSeeds", () => {
	it("takes the sample and statistic seeds from a result", () => {
		const result = {
			sample_seed: "1",
			statistics: {
				min_seed: "2",
				max_seed: "3",
				p25_seed: "4",
				p50_seed: "5",
				p75_seed: "6",
			},
		} as model.SimulationResult;
		expect(namedSeeds(result)).toEqual({
			sample: "1",
			min: "2",
			max: "3",
			p25: "4",
			p50: "5",
			p75: "6",
		});
	});
});
