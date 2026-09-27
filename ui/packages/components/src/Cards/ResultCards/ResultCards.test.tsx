import type { model } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (key: string) => key,
		i18n: { language: "en" },
	}),
}));

import TargetDPSCard from "./Damage/TargetDPSCard";
import RollupCards from "./RollupCards";

function result(dps: number): model.SimulationResult {
	return {
		statistics: {
			dps: { mean: dps },
			target_dps: { "1": { mean: dps } },
		},
	} as model.SimulationResult;
}

describe("result cards", () => {
	it("render the latest result after streaming updates", () => {
		const { rerender } = render(<RollupCards data={result(100)} />);
		rerender(<RollupCards data={result(200)} />);
		rerender(<RollupCards data={result(300)} />);
		expect(screen.getByText("300")).toBeInTheDocument();
		expect(screen.queryByText("100")).not.toBeInTheDocument();
	});

	it("drop the stale indicator once the deferred value catches up", () => {
		const { rerender } = render(<TargetDPSCard data={result(100)} />);
		rerender(<TargetDPSCard data={result(200)} />);
		expect(screen.getByText("result.dist")).toBeInTheDocument();
		expect(screen.queryByRole("status")).not.toBeInTheDocument();
	});
});
