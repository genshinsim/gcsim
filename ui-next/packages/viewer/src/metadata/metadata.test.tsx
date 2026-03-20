import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Commit, Iterations, Mode, Warnings } from "./metadata.js";

describe("Iterations", () => {
  it("renders iteration count", () => {
    render(<Iterations iterations={1000} />);
    expect(screen.getByTestId("iterations")).toHaveTextContent("1,000 iterations");
  });

  it("renders nothing when iterations is undefined", () => {
    const { container } = render(<Iterations />);
    expect(container.firstChild).toBeNull();
  });
});

describe("Mode", () => {
  it("renders SL for mode 0", () => {
    render(<Mode mode={0} />);
    expect(screen.getByTestId("mode")).toHaveTextContent("SL");
  });

  it("renders TTK for mode 1", () => {
    render(<Mode mode={1} />);
    expect(screen.getByTestId("mode")).toHaveTextContent("TTK");
  });

  it("renders nothing when mode is undefined", () => {
    const { container } = render(<Mode />);
    expect(container.firstChild).toBeNull();
  });
});

describe("Commit", () => {
  it("renders sim version and build date", () => {
    render(<Commit simVersion="2.5.0" buildDate="2026-03-19" />);
    expect(screen.getByTestId("sim-version")).toHaveTextContent("2.5.0");
    expect(screen.getByTestId("build-date")).toHaveTextContent("2026-03-19");
  });

  it("renders only sim version when build date is missing", () => {
    render(<Commit simVersion="2.5.0" />);
    expect(screen.getByTestId("sim-version")).toHaveTextContent("2.5.0");
    expect(screen.queryByTestId("build-date")).toBeNull();
  });

  it("renders nothing when both are missing", () => {
    const { container } = render(<Commit />);
    expect(container.firstChild).toBeNull();
  });
});

describe("Warnings", () => {
  it("renders active warnings as badges", () => {
    render(
      <Warnings
        warnings={{
          target_overlap: true,
          insufficient_energy: true,
          swap_cd: false,
        }}
      />,
    );
    expect(screen.getByText("Target Overlap")).toBeInTheDocument();
    expect(screen.getByText("Insufficient Energy")).toBeInTheDocument();
    expect(screen.queryByText("Swap CD")).not.toBeInTheDocument();
  });

  it("renders nothing when no warnings are active", () => {
    const { container } = render(
      <Warnings
        warnings={{
          target_overlap: false,
          insufficient_energy: false,
        }}
      />,
    );
    expect(container.querySelector("[data-testid='warnings']")).toBeNull();
  });

  it("renders nothing when warnings is undefined", () => {
    const { container } = render(<Warnings />);
    expect(container.firstChild).toBeNull();
  });
});
