import type { Sim } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { TeamHeader } from "./team-header.js";

const mockCharacters: Sim.Character[] = [
  {
    name: "hutao",
    level: 90,
    max_level: 90,
    element: "pyro",
    cons: 1,
    weapon: { name: "staffofhoma", refine: 1, level: 90, max_level: 90 },
    talents: { attack: 10, skill: 10, burst: 10 },
    stats: [],
    snapshot: [],
    sets: { crimsonwitchofflames: 4 },
  },
  {
    name: "xingqiu",
    level: 90,
    max_level: 90,
    element: "hydro",
    cons: 6,
    weapon: { name: "sacrificialsword", refine: 5, level: 90, max_level: 90 },
    talents: { attack: 1, skill: 10, burst: 13 },
    stats: [],
    snapshot: [],
    sets: { emblemofseveredfate: 4 },
  },
];

describe("TeamHeader", () => {
  it("renders a card for each character", () => {
    render(<TeamHeader characters={mockCharacters} />);
    const names = screen.getAllByTestId("char-name");
    expect(names).toHaveLength(2);
    expect(names[0]).toHaveTextContent("hutao");
    expect(names[1]).toHaveTextContent("xingqiu");
  });

  it("shows level, constellation, and weapon", () => {
    render(<TeamHeader characters={mockCharacters} />);
    const levels = screen.getAllByTestId("char-level");
    expect(levels[0]).toHaveTextContent("Lv. 90/90");

    const cons = screen.getAllByTestId("char-cons");
    expect(cons[0]).toHaveTextContent("C1");
    expect(cons[1]).toHaveTextContent("C6");

    const weapons = screen.getAllByTestId("char-weapon");
    expect(weapons[0]).toHaveTextContent("staffofhoma");

    const refines = screen.getAllByTestId("char-weapon-refine");
    expect(refines[0]).toHaveTextContent("1");
  });

  it("renders nothing when characters is undefined", () => {
    const { container } = render(<TeamHeader />);
    expect(container.firstChild).toBeNull();
  });

  it("renders nothing when characters is empty", () => {
    const { container } = render(<TeamHeader characters={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
