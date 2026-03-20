import { Card, CardContent, cn } from "@gcsim/primitives";
import type { Sim } from "@gcsim/types";

export interface DPSCardProps {
  characterName: string;
  stat?: Sim.FloatStat;
  /** Maximum DPS across all characters, used to calculate proportional bar width */
  maxDPS?: number;
  className?: string;
}

function formatDPS(value: number | undefined): string {
  if (value == null) return "—";
  return value.toLocaleString(undefined, { maximumFractionDigits: 0 });
}

export function DPSCard({ characterName, stat, maxDPS, className }: DPSCardProps) {
  const mean = stat?.mean;
  const barWidth = mean != null && maxDPS != null && maxDPS > 0 ? (mean / maxDPS) * 100 : 0;

  return (
    <Card className={cn("min-w-[200px]", className)} data-testid="dps-card">
      <CardContent className="p-3">
        <div className="flex items-center justify-between">
          <span className="font-medium capitalize" data-testid="dps-char-name">
            {characterName}
          </span>
          <span className="font-bold tabular-nums" data-testid="dps-value">
            {formatDPS(mean)}
          </span>
        </div>
        <div className="bg-muted mt-2 h-2 w-full overflow-hidden rounded-full">
          <div
            className="bg-primary h-full rounded-full transition-all"
            style={{ width: `${barWidth}%` }}
            data-testid="dps-bar"
          />
        </div>
      </CardContent>
    </Card>
  );
}
