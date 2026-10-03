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

	it("steps the strip a column with the arrow keys from anywhere on the page", async () => {
		renderLog();
		await userEvent.keyboard("{ArrowRight}");
		expect(scrolls).toHaveLength(1);
		expect(scrolls[0]).toBeGreaterThan(0);
		await userEvent.keyboard("{ArrowLeft}");
		expect(scrolls).toEqual([scrolls[0], 0]);
	});

	it("leaves arrow keys to the search box", async () => {
		renderLog();
		await userEvent.click(
			screen.getByPlaceholderText("sample.search_placeholder"),
		);
		await userEvent.keyboard("{ArrowRight}{ArrowLeft}");
		expect(scrolls).toEqual([]);
	});

	it("leaves arrow keys to an open frame view", async () => {
		renderLog();
		await userEvent.click(
			screen.getAllByRole("button", { name: /^sample\.open_frame / })[0],
		);
		await userEvent.keyboard("{ArrowRight}");
		// the frame view centres the strip on its new frame, and nothing else scrolls it
		expect(
			within(screen.getByRole("dialog")).getByText(/sample\.frame_title/),
		).toBeInTheDocument();
		expect(scrolls).toHaveLength(1);
	});

	it("restarts the match count when the filter changes the matches", async () => {
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		const next = screen.getByRole("button", { name: "sample.next_match" });
		await userEvent.click(next);
		await userEvent.click(next);
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.log_options" }),
		);
		await userEvent.click(
			within(screen.getByRole("dialog")).getByText("viewer.clear"),
		);
		expect(screen.getByText("0/0")).toBeInTheDocument();
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

	it("offers generate beside the filter only when a handler is given", async () => {
		const { unmount } = renderLog();
		expect(
			screen.queryByRole("button", { name: "viewer.generate" }),
		).toBeNull();
		unmount();

		const onGenerate = vi.fn();
		renderLog({ onGenerate });
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.generate" }),
		);
		expect(onGenerate).toHaveBeenCalledOnce();
	});

	it("offers only generate when there is no sample", async () => {
		const onGenerate = vi.fn();
		renderLog({ sample: null, onGenerate });
		expect(screen.queryByTitle("sample.sim_lane")).toBeNull();
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.generate" }),
		);
		expect(onGenerate).toHaveBeenCalledOnce();
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

	it("leads a damage chip with the compact amount and keeps the exact one in its title", async () => {
		Element.prototype.scrollTo = function (
			this: Element,
			opts?: ScrollToOptions | number,
		) {
			if (typeof opts === "object") {
				this.scrollLeft = opts.left ?? 0;
				this.dispatchEvent(new Event("scroll"));
			}
		} as typeof Element.prototype.scrollTo;
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		await userEvent.click(
			screen.getByRole("button", { name: "sample.next_match" }),
		);
		const chip = screen.getByTitle(
			"122 · damage: Sesshou Sakura Tick [8,815] (crit)",
		);
		expect(chip).toHaveTextContent(/8\.82K Sesshou Sakura Tick \(crit\)$/);
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

	it("keeps a duration highlight when Escape closes the details dialog", async () => {
		renderLog();
		const durations = screen.getAllByRole("button", {
			name: /^sample\.show_duration /,
		});
		await userEvent.click(durations[0]);
		const clear = () => screen.queryAllByLabelText(/^sample\.clear_duration /);
		expect(clear()).toHaveLength(1);

		await userEvent.click(chipsOfType("action")[0]);
		await userEvent.keyboard("{Escape}");
		expect(screen.queryByRole("dialog")).toBeNull();
		expect(clear()).toHaveLength(1);

		await userEvent.keyboard("{Escape}");
		expect(clear()).toHaveLength(0);
	});
});
