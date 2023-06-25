import { Card } from "@blueprintjs/core";
import { Character } from "@gcsim/types";
import ArtifactsIcon from "./ArtifactsIcon";

type Props = {
  chars?: Character[];
  handleLoaded: () => void;
};

export const Avatars = ({ chars, handleLoaded }: Props) => {
  return (
    <Card className="flex flex-row p-2 max-w-fit">
      {chars?.map((c) => <AvatarCard key={c.name} c={c} handleLoaded={handleLoaded} />)}
    </Card>
  );
};

type CardProps = {
  c: Character;
  handleLoaded: () => void;
};

const AvatarCard = ({ c, handleLoaded }: CardProps) => {
  const sets: string[] = [];
  let half = false;

  if ("sets" in c) {
    for (const [key,] of Object.entries(c.sets)) {
      sets.push(key);
    }
  }

  if (sets.length == 1 && c.sets[sets[0]] == 2) {
    half = true;
  }

  return (
    <div className="relative w-24">
      <img
        className="object-contain h-24"
        key={c.name}
        src={`/api/assets/avatar/${c.name}.png`}
        onLoad={handleLoaded}
      />
      <Card className={
        "absolute left-0 top-0 flex flex-row gap-1 px-1 py-0 rounded-lg " +
        "font-bold font-mono text-base "
      }>
        <span className="text-geo">{`C${c.cons}`}</span>
        <span className="text-electro">{`R${c.weapon.refine}`}</span>
      </Card>
      <div className="equip">
        <svg key={c.weapon.name} width="62" height="65">
          <filter id="outlinew">
            <feMorphology
              in="SourceAlpha"
              result="expanded"
              operator="dilate"
              radius="1"
            />
            <feFlood floodColor="white" />
            <feComposite in2="expanded" operator="in" />
            <feComposite in="SourceGraphic" />
          </filter>
          <filter id="outlineb">
            <feMorphology
              in="SourceAlpha"
              result="expanded"
              operator="dilate"
              radius="1.5"
            />
            <feFlood floodColor="black" />
            <feComposite in2="expanded" operator="in" />
            <feGaussianBlur stdDeviation="1" />
            <feComposite in="SourceGraphic" />
          </filter>
          <image
            filter="url(#outlinew) url(#outlineb)"
            href={`/api/assets/weapons/${c.weapon.name}.png`}
            height="58"
            width="58"
            x="3"
            y="3"
          />
          { sets.length > 0 ? <ArtifactsIcon sets={sets} half={half} /> : null}
        </svg>
      </div>
    </div>
  );
};