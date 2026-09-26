import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import { NameSearch } from "./NameSearch";

describe("NameSearch", () => {
	it("renders a trigger for each searchable category", () => {
		render(<NameSearch />);
		for (const label of [
			"db.characters",
			"simple.weapons",
			"simple.artifacts",
			"simple.enemies",
			"simple.actions",
			"simple.stats",
		]) {
			expect(screen.getByRole("button", { name: label })).toBeTruthy();
		}
	});
});
