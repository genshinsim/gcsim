import { Card, CardContent, CardHeader, CardTitle, cn } from "@gcsim/primitives";
import type { Sim } from "@gcsim/types";

export interface RollupCardProps {
  label: string;
  stat?: Sim.FloatStat | Sim.SummaryStat;
  className?: string;
}

function formatNumber(value: number | undefined): string {
  if (value == null) return "—";
  return value.toLocaleString(undefined, { maximumFractionDigits: 2 });
}

export function RollupCard({ label, stat, className }: RollupCardProps) {
  return (
    <Card className={cn("min-w-[180px]", className)} data-testid="rollup-card">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium">{label}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold" data-testid="rollup-mean">
          {formatNumber(stat?.mean)}
        </div>
        <div className="text-muted-foreground mt-1 text-xs">
          <span data-testid="rollup-min">Min: {formatNumber(stat?.min)}</span>
          {" / "}
          <span data-testid="rollup-max">Max: {formatNumber(stat?.max)}</span>
          {" / "}
          <span data-testid="rollup-sd">SD: {formatNumber(stat?.sd)}</span>
        </div>
      </CardContent>
    </Card>
  );
}
