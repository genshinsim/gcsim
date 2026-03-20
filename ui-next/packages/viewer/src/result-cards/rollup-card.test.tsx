import type { Sim } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { RollupCard } from "./rollup-card.js";

const mockStat: Sim.SummaryStat = {
  min: 35000,
  max: 65000,
  mean: 50250.5,
  sd: 4200.3,
};

describe("RollupCard", () => {
  it("renders the label", () => {
    render(<RollupCard label="DPS" stat={mockStat} />);
    expect(screen.getByText("DPS")).toBeInTheDocument();
  });

  it("renders mean value", () => {
    render(<RollupCard label="DPS" stat={mockStat} />);
    expect(screen.getByTestId("rollup-mean")).toHaveTextContent("50,250.5");
  });

  it("renders min, max, and sd", () => {
    render(<RollupCard label="DPS" stat={mockStat} />);
    expect(screen.getByTestId("rollup-min")).toHaveTextContent("Min: 35,000");
    expect(screen.getByTestId("rollup-max")).toHaveTextContent("Max: 65,000");
    expect(screen.getByTestId("rollup-sd")).toHaveTextContent("SD: 4,200.3");
  });

  it("renders dashes when stat is undefined", () => {
    render(<RollupCard label="DPS" />);
    expect(screen.getByTestId("rollup-mean")).toHaveTextContent("—");
    expect(screen.getByTestId("rollup-min")).toHaveTextContent("Min: —");
  });
});
