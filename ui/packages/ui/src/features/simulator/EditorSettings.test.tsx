// @vitest-environment jsdom
import "@gcsim/components/vitest.setup";
import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (k: string, o?: { src?: string }) => (o?.src ? `${k}:${o.src}` : k),
		i18n: { language: "en" },
	}),
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

async function openImportMenu() {
	render(<EditorSettings />);
	await userEvent.click(screen.getByRole("button", { name: "simple.import" }));
}

afterEach(cleanup);

describe("EditorSettings", () => {
	it("shows only the Import button", () => {
		render(<EditorSettings />);
		expect(screen.getAllByRole("button").map((b) => b.textContent)).toEqual([
			"simple.import",
		]);
	});

	it("lists only the GO and Enka imports in the Import dropdown", async () => {
		await openImportMenu();
		expect(screen.getAllByRole("menuitem").map((m) => m.textContent)).toEqual([
			"simple.tools_import:GO",
			"simple.tools_import:Enka",
		]);
		expect(screen.queryByRole("separator")).toBeNull();
	});

	it.each([
		["GO", "good-dialog"],
		["Enka", "enka-dialog"],
	])("opens the %s import dialog from the dropdown", async (src, testId) => {
		await openImportMenu();
		await userEvent.click(
			screen.getByRole("menuitem", { name: `simple.tools_import:${src}` }),
		);
		expect(screen.getByTestId(testId)).toBeInTheDocument();
	});
});
