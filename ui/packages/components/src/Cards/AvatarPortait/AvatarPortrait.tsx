import {model} from '@gcsim/types';
import {DataColorsConst} from '../../common/gcsim';
import {Badge} from '../../common/ui/badge';
import {charBG} from '../../lib/helper';
import ArtifactsIcon from './ArtifactsIcon';

const WeaponImage = ({weapon}: {weapon: model.Weapon}) => {
  return (
    <div className="absolute bottom-[-4px] w-[62] right-[-12px] opacity-85">
      <svg
        key={weapon.name}
        width="62"
        height="65"
        onError={(e: React.SyntheticEvent<SVGSVGElement, Event>) => {
          (e.target as SVGAElement).href.baseVal =
            '/api/assets/misc/default.png';
        }}>
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
          href={`/api/assets/weapons/${weapon.name}.png`}
          height="55"
          width="55"
          x="0"
          y="3"
        />
      </svg>
    </div>
  );
};

type AvatarPortraitProps = {
  char: model.Character | null;
  i: number;
  invalid: boolean;
  onImageLoaded: () => void;

  //optional classes
  className?: string;
};

export const AvatarPortrait = ({
  char,
  i,
  invalid,
  onImageLoaded,
  className = '',
}: AvatarPortraitProps) => {
  //display an empty card here
  if (char === null) {
    return (
      <div
        className={
          'flex flex-col bg-gray-400 border border-gray-600 rounded-sm' +
          (className == '' ? '' : ' ' + className)
        }>
        <div className="flex justify-center">
          <img
            src={'/api/assets/misc/nahida.png'}
            className=" object-contain opacity-50 h-24"
            onLoad={onImageLoaded}
          />
        </div>
      </div>
    );
  }
  const sets: string[] = [];
  let half = false;

  if (char.sets && char.sets !== null) {
    for (const [key] of Object.entries(char.sets)) {
      sets.push(key);
    }
    if (sets.length == 1 && char.sets[sets[0]] == 2) {
      half = true;
    }
  }

  return (
    <div
      className={
        'flex flex-col bg-gray-400 border border-gray-600 rounded-sm' +
        (className == '' ? '' : ' ' + className)
      }>
      <div className={`relative w-full pt-2 z-0 ${charBG(char.element ?? '')}`}>
        <div
          className="absolute top-0 left-0 right-0 bottom-0 !bg-cover !bg-center mix-blend-luminosity z-1 "
          style={{background: `url(/api/assets/misc/overlay.jpg)`}}></div>

        <div className="flex justify-center">
          <img
            className="relative object-contain h-24"
            key={char.name}
            src={`/api/assets/avatar/${char.name}.png`}
            onError={(e) => {
              (e.target as HTMLImageElement).src =
                '/api/assets/misc/default.png';
              onImageLoaded();
            }}
            onLoad={onImageLoaded}
          />
        </div>
        {char.weapon ? <WeaponImage weapon={char.weapon} /> : null}
        <div className="absolute bottom-0 left-0 opacity-85">
          <svg
            width={35}
            height={35}
            onError={(e: React.SyntheticEvent<SVGSVGElement, Event>) => {
              (e.target as SVGAElement).href.baseVal =
                '/api/assets/misc/default.png';
            }}>
            {sets.length > 0 ? <ArtifactsIcon sets={sets} half={half} /> : null}
          </svg>
        </div>

        <div
          className={
            'absolute left-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none ' +
            'font-bold font-mono text-xs rounded-tl-sm rounded-br-lg ' +
            'bg-gray-700 opacity-85'
          }>
          <div className="flex flex-row gap-1 min-h-fit">
            <span className="text-geo">{`C${char.cons ?? 0}`}</span>
            {char.weapon ? (
              <span className="text-electro">{`R${
                char.weapon.refine ?? 0
              }`}</span>
            ) : null}
          </div>
        </div>

        <div
          className={
            'absolute right-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none ' +
            'font-mono text-xs rounded-tr-sm rounded-bl-lg ' +
            'bg-gray-700 opacity-85'
          }>
          <div className="flex flex-row gap-1 min-h-fit items-center">
            <span className="text-xs text-gray-400">lvl</span>
            <span
              className={`font-bold`}
              style={{color: DataColorsConst.qualitative5(i)}}>
              {char.level}
            </span>
          </div>
        </div>

        {invalid && (
          <div className="absolute left-0 top-1/3 w-full">
            <Badge className="flex flex-row items-center justify-center gap-2">
              <div className="font-mono select-none text-red-500 font-bold text-xs uppercase">
                WIP
              </div>
            </Badge>
          </div>
        )}
      </div>
    </div>
  );
};
