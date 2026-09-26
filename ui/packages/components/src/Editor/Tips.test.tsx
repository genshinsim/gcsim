import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import { ActionListTip, TeamTip } from "./Tips";

describe("ActionListTip", () => {
	it("links to Discord and the docs", () => {
		render(<ActionListTip />);
		expect(screen.getByText("Discord").getAttribute("href")).toBe(
			"https://discord.gg/W36ZwwhEaG",
		);
		expect(screen.getByText("simple.documentation").getAttribute("href")).toBe(
			"https://docs.gcsim.app/guides",
		);
	});

	it("calls onHide from the hide button when provided", async () => {
		const onHide = vi.fn();
		render(<ActionListTip onHide={onHide} />);
		await userEvent.click(screen.getByText("simple.hide_all_tips"));
		expect(onHide).toHaveBeenCalledTimes(1);
	});

	it("omits the hide button without an onHide handler", () => {
		render(<ActionListTip />);
		expect(screen.queryByText("simple.hide_all_tips")).toBeNull();
	});
});

describe("TeamTip", () => {
	it("shows the team builder blurb", () => {
		render(<TeamTip />);
		expect(screen.getByText("simple.video")).toBeTruthy();
	});

	it("calls onHide from the hide button when provided", async () => {
		const onHide = vi.fn();
		render(<TeamTip onHide={onHide} />);
		await userEvent.click(screen.getByText("simple.hide_all_tips"));
		expect(onHide).toHaveBeenCalledTimes(1);
	});
});
