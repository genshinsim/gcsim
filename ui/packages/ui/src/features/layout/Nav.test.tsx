// @vitest-environment jsdom
import "@gcsim/components/vitest.setup";
import { cleanup, render, screen } from "@testing-library/react";
import type { ComponentProps } from "react";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (k: string) => k,
		i18n: { language: "en" },
	}),
}));
vi.mock("@tanstack/react-router", () => ({
	Link: ({ to, children, ...props }: ComponentProps<"a"> & { to: string }) => (
		<a href={to} {...props}>
			{children}
		</a>
	),
}));
vi.mock("./ExecutorDialog", () => ({
	ExecutorDialog: () => <button type="button">executor</button>,
}));

import Nav from "./Nav";

afterEach(cleanup);

describe("Nav", () => {
	it("shows Sample Upload as an icon link between the executor and settings", () => {
		render(<Nav />);
		const upload = screen.getByRole("link", {
			name: "simple.tools_sample_upload",
		});
		expect(upload).toHaveAttribute("href", "/sample/upload");
		expect(upload.textContent).toBe("");

		const group = upload.parentElement as HTMLElement;
		const controls = Array.from(group.querySelectorAll("a, button")).map(
			(el) => el.getAttribute("aria-label") ?? el.textContent,
		);
		expect(controls.slice(0, 3)).toEqual([
			"executor",
			"simple.tools_sample_upload",
			"simple.settings",
		]);
	});
});
