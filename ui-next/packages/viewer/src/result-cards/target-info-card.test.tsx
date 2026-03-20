import type { Sim } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { TargetInfoCard } from "./target-info-card.js";

const mockEnemies: Sim.Enemy[] = [
  {
    name: "target-1",
    level: 100,
    hp: 10000000,
    resist: {
      pyro: 0.1,
      hydro: 0.1,
      cryo: 0.1,
    },
    position: { x: 0, y: 0, r: 1 },
  },
];

describe("TargetInfoCard", () => {
  it("renders target name", () => {
    render(<TargetInfoCard enemies={mockEnemies} />);
    expect(screen.getByTestId("target-name")).toHaveTextContent("target-1");
  });

  it("renders target level", () => {
    render(<TargetInfoCard enemies={mockEnemies} />);
    expect(screen.getByTestId("target-level")).toHaveTextContent("Level 100");
  });

  it("renders resistances", () => {
    render(<TargetInfoCard enemies={mockEnemies} />);
    const resists = screen.getByTestId("target-resists");
    expect(resists).toHaveTextContent("pyro: 10%");
    expect(resists).toHaveTextContent("hydro: 10%");
    expect(resists).toHaveTextContent("cryo: 10%");
  });

  it("uses fallback name when name is not provided", () => {
    render(<TargetInfoCard enemies={[{ level: 90, resist: {} }]} />);
    expect(screen.getByTestId("target-name")).toHaveTextContent("Target 1");
  });

  it("renders nothing when enemies is undefined", () => {
    const { container } = render(<TargetInfoCard />);
    expect(container.firstChild).toBeNull();
  });

  it("renders nothing when enemies is empty", () => {
    const { container } = render(<TargetInfoCard enemies={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("renders multiple targets", () => {
    const twoEnemies: Sim.Enemy[] = [
      { name: "boss-1", level: 100, resist: {} },
      { name: "boss-2", level: 95, resist: {} },
    ];
    render(<TargetInfoCard enemies={twoEnemies} />);
    const names = screen.getAllByTestId("target-name");
    expect(names).toHaveLength(2);
    expect(names[0]).toHaveTextContent("boss-1");
    expect(names[1]).toHaveTextContent("boss-2");
  });
});
