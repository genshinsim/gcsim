import { Character } from '~src/Types/sim';
import { useTranslation } from 'react-i18next';

type CharacterTooltipProps = {
  char: Character;
};

export function CharacterTooltip({ char }: CharacterTooltipProps) {
  let { t } = useTranslation();
  return (
    <div className="m-2 flex flex-col">
      <div className="ml-auto font-bold capitalize">{`${t(
        'game:character_names.' + char.name
      )} ${t('db.c_pre')}${char.cons}${t('db.c_post')} ${char.talents.attack}/${
        char.talents.skill
      }/${char.talents.burst}`}</div>
      <div className="w-full border-b border-gray-500 mt-2 mb-2"></div>
      <div className="capitalize flex flex-row">
        <img
          src={'/images/weapons/' + char.weapon.name + '.png'}
          alt={char.weapon.name}
          className="wide:h-8 h-auto "
        />
        <div className="mt-auto mb-auto">
          {t('game:weapon_names.' + char.weapon.name) +
            t('db.r') +
            char.weapon.refine}
        </div>
      </div>
      {/* <div className="ml-auto">{`${t('db.er')}${char.er * 100 + 100}%`}</div> */}
    </div>
  );
}
