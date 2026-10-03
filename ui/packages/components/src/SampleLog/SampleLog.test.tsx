import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (k: string, o?: Record<string, unknown>) =>
			o == null ? k : `${k} ${JSON.stringify(o)}`,
	}),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import { SampleLog } from "./SampleLog";
import { loadSampleFilter } from "./sampleFilter";
import { sampleFixture as sample } from "./testdata";

const scrolls: number[] = [];

beforeEach(() => {
	localStorage.clear();
	scrolls.length = 0;
	Element.prototype.scrollTo = ((opts?: ScrollToOptions | number) => {
		if (typeof opts === "object") {
			scrolls.push(opts.left ?? 0);
		}
	}) as typeof Element.prototype.scrollTo;
});

function renderLog(
	props: Partial<React.ComponentProps<typeof SampleLog>> = {},
) {
	return render(<SampleLog sample={sample} {...props} />);
}

const chipsOfType = (type: string) =>
	screen.queryAllByTitle(new RegExp(`^\\d+ · ${type}: `));

describe("SampleLog", () => {
	it("draws a lane per character plus the sim lane", () => {
		renderLog();
		for (const c of sample.character_details ?? []) {
			expect(
				screen.getByTitle(`game:character_names.${c.name}`),
			).toBeInTheDocument();
		}
		expect(screen.getAllByTitle("sample.sim_lane").length).toBeGreaterThan(0);
	});

	it("steps through search matches, centring the strip on each", async () => {
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		const next = screen.getByRole("button", { name: "sample.next_match" });
		await userEvent.click(next);
		await userEvent.click(next);
		expect(scrolls).toHaveLength(2);
		expect(screen.getByText(/^2\/\d+$/)).toBeInTheDocument();
	});

	it("offers download only when a handler is given", async () => {
		const { unmount } = renderLog();
		expect(
			screen.queryByRole("button", { name: "viewer.download" }),
		).toBeNull();
		unmount();

		const onDownload = vi.fn();
		renderLog({ onDownload });
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.download" }),
		);
		expect(onDownload).toHaveBeenCalledExactlyOnceWith(sample);
	});

	it("saves the filter and hides the types it leaves out", async () => {
		renderLog();
		expect(chipsOfType("action").length).toBeGreaterThan(0);
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.log_options" }),
		);
		const popover = screen.getByRole("dialog");
		await userEvent.click(within(popover).getByText("action"));
		expect(loadSampleFilter(localStorage)).not.toContain("action");
		expect(chipsOfType("action")).toEqual([]);
	});

	it("opens an event's details from its chip", async () => {
		renderLog();
		const chip = chipsOfType("action")[0];
		await userEvent.click(chip);
		const heading = chip.title.replace(" · action: ", " · ");
		expect(
			within(screen.getByRole("dialog")).getByRole("heading", {
				name: heading,
			}),
		).toBeInTheDocument();
	});
});
