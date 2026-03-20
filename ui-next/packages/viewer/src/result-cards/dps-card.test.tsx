import type { Sim } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { DPSCard } from "./dps-card.js";

const mockDPS: Sim.FloatStat = {
  min: 20000,
  max: 45000,
  mean: 35100,
  sd: 3200,
};

describe("DPSCard", () => {
  it("renders character name", () => {
    render(<DPSCard characterName="hutao" stat={mockDPS} maxDPS={35100} />);
    expect(screen.getByTestId("dps-char-name")).toHaveTextContent("hutao");
  });

  it("renders formatted DPS value", () => {
    render(<DPSCard characterName="hutao" stat={mockDPS} maxDPS={35100} />);
    expect(screen.getByTestId("dps-value")).toHaveTextContent("35,100");
  });

  it("renders proportional bar width", () => {
    render(<DPSCard characterName="hutao" stat={mockDPS} maxDPS={70200} />);
    const bar = screen.getByTestId("dps-bar");
    expect(bar).toHaveStyle({ width: "50%" });
  });

  it("renders 100% bar when DPS equals max", () => {
    render(<DPSCard characterName="hutao" stat={mockDPS} maxDPS={35100} />);
    const bar = screen.getByTestId("dps-bar");
    expect(bar).toHaveStyle({ width: "100%" });
  });

  it("renders dash when stat is undefined", () => {
    render(<DPSCard characterName="hutao" />);
    expect(screen.getByTestId("dps-value")).toHaveTextContent("—");
  });

  it("renders 0% bar when maxDPS is undefined", () => {
    render(<DPSCard characterName="hutao" stat={mockDPS} />);
    const bar = screen.getByTestId("dps-bar");
    expect(bar).toHaveStyle({ width: "0%" });
  });
});
