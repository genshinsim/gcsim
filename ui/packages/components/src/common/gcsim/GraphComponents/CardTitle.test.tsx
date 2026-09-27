import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (key: string) => key }),
}));

import CardTitle from "./CardTitle";

describe("CardTitle", () => {
	it("shows the updating indicator while stale", () => {
		render(<CardTitle title="DPS" stale />);
		expect(screen.getByText("DPS")).toBeInTheDocument();
		expect(screen.getByRole("status")).toHaveAccessibleName("result.updating");
	});

	it("hides the indicator once fresh", () => {
		render(<CardTitle title="DPS" stale={false} />);
		expect(screen.getByText("DPS")).toBeInTheDocument();
		expect(screen.queryByRole("status")).not.toBeInTheDocument();
	});
});
