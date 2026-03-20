import { Card, CardContent, CardHeader, CardTitle, cn } from "@gcsim/primitives";
import type { Sim } from "@gcsim/types";

export interface TargetInfoCardProps {
  enemies?: Sim.Enemy[];
  className?: string;
}

function formatResist(value: number): string {
  return `${(value * 100).toFixed(0)}%`;
}

export function TargetInfoCard({ enemies, className }: TargetInfoCardProps) {
  if (!enemies || enemies.length === 0) return null;

  return (
    <div data-testid="target-info" className={cn("flex flex-wrap gap-2", className)}>
      {enemies.map((enemy, i) => {
        const label = enemy.name ?? `Target ${i + 1}`;
        return (
          <Card key={label} className="min-w-[180px] flex-1">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium" data-testid="target-name">
                {label}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {enemy.level != null && (
                <div className="text-muted-foreground text-sm" data-testid="target-level">
                  Level {enemy.level}
                </div>
              )}
              {enemy.resist && Object.keys(enemy.resist).length > 0 && (
                <div
                  className="mt-1 flex flex-wrap gap-x-3 gap-y-0.5 text-xs"
                  data-testid="target-resists"
                >
                  {Object.entries(enemy.resist).map(([element, value]) => (
                    <span key={element} className="capitalize">
                      {element}: {formatResist(value)}
                    </span>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
