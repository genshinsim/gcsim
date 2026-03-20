import { Badge } from "@gcsim/primitives";
import type { Sim } from "@gcsim/types";

export interface IterationsProps {
  iterations?: number;
}

export function Iterations({ iterations }: IterationsProps) {
  if (iterations == null) return null;
  return <span data-testid="iterations">{iterations.toLocaleString()} iterations</span>;
}

export interface ModeProps {
  mode?: number;
}

export function Mode({ mode }: ModeProps) {
  if (mode == null) return null;
  const label = mode === 0 ? "SL" : "TTK";
  return <span data-testid="mode">{label}</span>;
}

export interface CommitProps {
  simVersion?: string;
  buildDate?: string;
}

export function Commit({ simVersion, buildDate }: CommitProps) {
  if (!simVersion && !buildDate) return null;
  return (
    <span data-testid="commit">
      {simVersion && <span data-testid="sim-version">{simVersion}</span>}
      {simVersion && buildDate && " — "}
      {buildDate && <span data-testid="build-date">{buildDate}</span>}
    </span>
  );
}

export interface WarningsProps {
  warnings?: Sim.Warnings;
}

const warningLabels: Record<string, string> = {
  target_overlap: "Target Overlap",
  insufficient_energy: "Insufficient Energy",
  insufficient_stamina: "Insufficient Stamina",
  swap_cd: "Swap CD",
  skill_cd: "Skill CD",
  dash_cd: "Dash CD",
  burst_cd: "Burst CD",
};

export function Warnings({ warnings }: WarningsProps) {
  if (!warnings) return null;

  const active = Object.entries(warnings).filter(([, v]) => v === true);
  if (active.length === 0) return null;

  return (
    <div data-testid="warnings" className="flex flex-wrap gap-1">
      {active.map(([key]) => (
        <Badge key={key} variant="destructive">
          {warningLabels[key] ?? key}
        </Badge>
      ))}
    </div>
  );
}
