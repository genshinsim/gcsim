import { Card, Tag } from "@blueprintjs/core";
import { Character } from "@gcsim/types";
import { DataColorsConst } from "@gcsim/ui/src/Pages/Viewer/Components/Util";
import ArtifactsIcon from "./ArtifactsIcon";
import placeholder from "./default.png";
import overlay from "./overlay.jpg";

type Props = {
  chars?: Character[];
  invalid?: string[];
  handleLoaded: () => void;
};

export const Avatars = ({ chars, invalid, handleLoaded }: Props) => {
  return (
    <div className="flex flex-row grow gap-2 justify-between">
      {chars?.map((c, i) => {
        return (
          <AvatarCard
              key={c.name}
              c={c}
              i={i}
              invalid={invalid != null && invalid?.includes(c.name)}
              handleLoaded={handleLoaded} />
        );
      })}
    </div>
  );
};

type CardProps = {
  c: Character;
  i: number;
  invalid: boolean;
  handleLoaded: () => void;
};

const AvatarCard = ({ c, i, invalid, handleLoaded }: CardProps) => {
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
    <div className="flex flex-col w-full bg-bp4-dark-gray-400 shadow border border-gray-600 rounded-sm">
      <div className={`relative w-full pt-2 z-0 ${charBG(c.element)}`}>
        <div
            className="absolute top-0 left-0 right-0 bottom-0 !bg-cover !bg-center mix-blend-luminosity z-1 "
            style={{ background: `url(${overlay})` }}>
        </div>

        <div className="flex justify-center">
          <img
            className="relative object-contain h-24"
            key={c.name}
            src={`/api/assets/avatar/${c.name}.png`}
            onError={(e) => {
              (e.target as HTMLImageElement).src = placeholder;
              handleLoaded();
            }}
            onLoad={handleLoaded}
          />
        </div>
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
              x="0"
              y="3"
            />
          </svg>
        </div>
        <div className="absolute bottom-0 left-0">
          <svg width={35} height={35}>
            { sets.length > 0 ? <ArtifactsIcon sets={sets} half={half} /> : null}
          </svg>
        </div>

        <Card className={
          "absolute left-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none " +
          "font-bold font-mono text-base rounded-tl-sm rounded-br-lg "
        }>
          <div className="flex flex-row gap-1 min-h-fit">
            <span className="text-geo">{`C${c.cons ?? 0}`}</span>
            <span className="text-electro">{`R${c.weapon.refine ?? 0}`}</span>
          </div>
        </Card>

        <Card className={
          "absolute right-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none " +
          "font-mono text-base rounded-tr-sm rounded-bl-lg "
        }>
          <div className="flex flex-row gap-1 min-h-fit items-center">
            <span className="text-xs text-gray-400">lvl</span>
            <span className={`font-bold`} style={{ color: DataColorsConst.qualitative5(i) }}>
              {c.level}
            </span>
          </div>
        </Card>

      {invalid &&
        <div className="absolute left-0 top-1/3 w-full">
          <Tag large={true} intent="danger" fill>
            <div className="flex flex-row items-center justify-center gap-2 font-mono select-none">
              <div className="font-bold text-sm uppercase">In Progress</div>
            </div>
          </Tag>
        </div>
      }
      </div>
    </div>
  );
};

function charBG(element: string) {
  switch (element) {
    case "cryo":
      return "bg-gradient-to-r from-gray-700 to-blue-300";
    case "hydro":
      return "bg-gradient-to-r from-gray-700 to-blue-500";
    case "pyro":
      return "bg-gradient-to-r from-gray-700 to-red-400";
    case "electro":
      return "bg-gradient-to-r from-gray-700 to-purple-300";
    case "anemo":
      return "bg-gradient-to-r from-gray-700 to-teal-500";
    case "dendro":
      return "bg-gradient-to-r from-gray-700 to-lime-700";
    case "geo":
      return "bg-gradient-to-r from-gray-700 to-yellow-400";
  }
  return "bg-gray-700";
}