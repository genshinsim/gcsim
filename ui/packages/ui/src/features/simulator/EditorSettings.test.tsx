// @vitest-environment jsdom
import "@gcsim/components/vitest.setup";
import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

const navigate = vi.fn();

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k, i18n: { language: "en" } }),
}));
vi.mock("@tanstack/react-router", () => ({
	useNavigate: () => navigate,
}));
vi.mock("./good/ImportFromGOODDialog", () => ({
	ImportFromGOODDialog: ({ isOpen }: { isOpen: boolean }) =>
		isOpen ? <div data-testid="good-dialog" /> : null,
}));
vi.mock("./enka/ImportFromEnkaDialog", () => ({
	ImportFromEnkaDialog: ({ isOpen }: { isOpen: boolean }) =>
		isOpen ? <div data-testid="enka-dialog" /> : null,
}));

import { EditorSettings } from "./EditorSettings";

afterEach(() => {
	cleanup();
	navigate.mockReset();
});

describe("EditorSettings", () => {
	it("shows Sample Upload as its own button left of Import", async () => {
		render(<EditorSettings />);
		const buttons = screen.getAllByRole("button");
		expect(buttons.map((b) => b.textContent)).toEqual([
			"simple.tools_sample_upload",
			"simple.import",
		]);

		await userEvent.click(buttons[0]);
		expect(navigate).toHaveBeenCalledWith({ to: "/sample/upload" });
	});

	it("lists only the GO and Enka imports in the Import dropdown", async () => {
		render(<EditorSettings />);
		await userEvent.click(
			screen.getByRole("button", { name: "simple.import" }),
		);

		expect(screen.getAllByRole("menuitem")).toHaveLength(2);
		expect(screen.queryByRole("separator")).toBeNull();

		await userEvent.click(screen.getAllByRole("menuitem")[0]);
		expect(screen.getByTestId("good-dialog")).toBeInTheDocument();
	});

	it("opens the Enka dialog from the dropdown", async () => {
		render(<EditorSettings />);
		await userEvent.click(
			screen.getByRole("button", { name: "simple.import" }),
		);
		await userEvent.click(screen.getAllByRole("menuitem")[1]);
		expect(screen.getByTestId("enka-dialog")).toBeInTheDocument();
	});
});
