import { Card, CardContent, cn } from "@gcsim/primitives";
import type { Sim } from "@gcsim/types";

export interface TeamHeaderProps {
  characters?: Sim.Character[];
  className?: string;
}

export function TeamHeader({ characters, className }: TeamHeaderProps) {
  if (!characters || characters.length === 0) return null;

  return (
    <div data-testid="team-header" className={cn("flex flex-wrap gap-2", className)}>
      {characters.map((char) => (
        <Card key={char.name} className="min-w-[140px] flex-1">
          <CardContent className="p-3">
            <div className="font-semibold capitalize" data-testid="char-name">
              {char.name}
            </div>
            <div className="text-muted-foreground text-sm">
              <span data-testid="char-level">
                Lv. {char.level}/{char.max_level}
              </span>
              {" — "}
              <span data-testid="char-cons">C{char.cons}</span>
            </div>
            <div className="text-muted-foreground text-sm">
              <span data-testid="char-weapon" className="capitalize">
                {char.weapon.name}
              </span>
              {" R"}
              <span data-testid="char-weapon-refine">{char.weapon.refine}</span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
