import ArtifactsIcon from "./ArtifactsIcon";

type CardProps = {
  c: any; //TODO(kyle): This should be typed model.ICharacter
  handleLoaded: () => void;
};

const AvatarCard = ({ c, handleLoaded }: CardProps) => {
  let half = false;
  let sets: string[] = [];

  for (const [key, value] of Object.entries(c.sets)) {
    sets.push(key);
  }

  if (sets.length == 1 && c.sets[sets[0]] == 2) {
    half = true;
  }

  return (
    <div className="card">
      <div className="char">
        <img
          key={c.name}
          src={`/api/assets/avatar/${c.name}.png`}
          onLoad={handleLoaded}
        />
      </div>
      <div className="equip">
        <svg key={c.weapon.name} width="91" height="95">
          <filter id="outlinew">
            <feMorphology
              in="SourceAlpha"
              result="expanded"
              operator="dilate"
              radius="1"
            />
            <feFlood flood-color="white" />
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
            <feFlood flood-color="black" />
            <feComposite in2="expanded" operator="in" />
            <feGaussianBlur stdDeviation="1" />
            <feComposite in="SourceGraphic" />
          </filter>
          <image
            filter="url(#outlinew) url(#outlineb)"
            href={`/api/assets/weapons/${c.weapon.name}.png`}
            height="85"
            width="85"
            x="3"
            y="3"
          />
          {<ArtifactsIcon sets={sets} half={half} />}
        </svg>
      </div>
    </div>
  );
};

export default AvatarCard;
