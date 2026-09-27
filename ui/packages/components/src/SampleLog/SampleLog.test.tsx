import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

const scrolls: { id: string; index: number }[] = [];
vi.mock("@tanstack/react-virtual", async () => {
	const { useId } = await import("react");
	return {
		useVirtualizer: () => {
			const id = useId();
			return {
				scrollToIndex: (index: number) => scrolls.push({ id, index }),
				getVirtualItems: () => [],
				getTotalSize: () => 0,
				measureElement: () => {},
			};
		},
	};
});

import { SampleLog } from "./SampleLog";
import { sampleFixture as sample } from "./testdata";

function renderLog(
	props: Partial<React.ComponentProps<typeof SampleLog>> = {},
) {
	return render(
		<SampleLog
			sample={sample}
			settings={["damage"]}
			onSettingsChange={() => {}}
			{...props}
		/>,
	);
}

async function search(root: HTMLElement, text: string) {
	const input = within(root).getByRole("textbox");
	await userEvent.clear(input);
	await userEvent.type(input, text);
	await userEvent.click(
		within(root).getByRole("button", { name: "search next" }),
	);
}

beforeEach(() => {
	scrolls.length = 0;
});

describe("SampleLog", () => {
	it("moves to the next match on each search", async () => {
		const { container } = renderLog();
		await search(container, "Sesshou");
		await search(container, "Sesshou");
		expect(scrolls).toHaveLength(2);
		expect(scrolls[1].index).toBeGreaterThan(scrolls[0].index);
	});

	it("keeps each instance's search position separate", async () => {
		const a = renderLog();
		const b = renderLog();
		await search(a.container, "Sesshou");
		await search(a.container, "Sesshou");
		await search(b.container, "Sesshou");

		const [a1, a2, b1] = scrolls;
		expect(a1.id).not.toBe(b1.id);
		expect(a2.index).toBeGreaterThan(a1.index);
		expect(b1.index).toBe(a1.index);
	});

	it("resets the search position", async () => {
		const { container } = renderLog();
		await search(container, "Sesshou");
		await userEvent.click(
			within(container).getByRole("button", { name: "reset search" }),
		);
		await search(container, "Sesshou");
		expect(scrolls.map((s) => s.index)).toEqual([
			scrolls[0].index,
			0,
			scrolls[0].index,
		]);
	});

	it("offers download only when a handler is given", async () => {
		const onDownload = vi.fn();
		const { unmount } = renderLog();
		expect(screen.queryByText("viewer.download")).toBeNull();
		unmount();

		renderLog({ onDownload });
		await userEvent.click(screen.getByText("viewer.download"));
		expect(onDownload).toHaveBeenCalledOnce();
	});

	it("reports preset changes through onSettingsChange", async () => {
		const onSettingsChange = vi.fn();
		renderLog({ onSettingsChange });
		await userEvent.click(screen.getByText("simple.settings"));
		await userEvent.click(screen.getByText("viewer.clear"));
		expect(onSettingsChange).toHaveBeenCalledWith([]);
	});
});
