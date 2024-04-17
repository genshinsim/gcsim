import {
  Button,
  Collapse,
  Drawer,
  DrawerSize,
  Intent,
  Position,
} from '@blueprintjs/core';
import tagData from '@gcsim/data/src/tags.json';
import {useContext, useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {FaArrowDown, FaArrowUp, FaFilter, FaSearch} from 'react-icons/fa';
import useDebounce from '../SharedHooks/debounce';

import {
  charNames,
  FilterContext,
  FilterDispatchContext,
  ItemFilterState,
  SortByDirection,
  sortByParams,
} from './FilterComponents/Filter.utils';

export function Filter() {
  // https://github.com/i18next/next-i18next/issues/1795
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);

  const dispatch = useContext(FilterDispatchContext);
  const [isOpen, setIsOpen] = useState(false);

  // const filter = useContext(FilterContext);
  // const includedCharacterFilters: CharFilterState[] = Object.keys(
  //   filter.charFilter
  // )
  //   .filter((key) => filter.charFilter[key].state === ItemFilterState.include)
  //   .map((key) => filter.charFilter[key]);

  const [value, setValue] = useState<string>('');
  const debouncedValue = useDebounce<string>(value, 500);

  useEffect(() => {
    dispatch({type: 'setCustomFilter', customFilter: debouncedValue});
  }, [debouncedValue, dispatch]);

  return (
    <div>
      <button
        className="flex flex-row gap-2 bp4-button justify-center items-center p-3 bp4-intent-primary h-12 w-12"
        onClick={() => setIsOpen(!isOpen)}>
        <FaFilter size={24} className="opacity-80" />
      </button>

      <Drawer
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose
        autoFocus
        isCloseButtonShown
        title={
          <div
            className="flex flex-row justify-between
          ">
            <div className="text-xl pb-1 ">{t('db.filter')}</div>
            <ClearFilterButton />
          </div>
        }
        onClose={() => setIsOpen(false)}
        position={Position.LEFT}
        size={DrawerSize.SMALL}>
        <div className="flex flex-col gap-2 overflow-y-auto overflow-x-hidden p-2">
          <input
            className="bp4-input bp4-icon bp4-icon-filter"
            placeholder={t('db.customFilter')}
            type="text"
            dir="auto"
            onChange={(e) => {
              setValue(e.target.value);
            }}
          />
          <CharacterFilter />
          <TagFilter />
          <SortBy />
        </div>
      </Drawer>
    </div>
  );
}

function ClearFilterButton() {
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);
  const dispatch = useContext(FilterDispatchContext);
  return (
    <button
      className="bp4-button bp4-intent-danger bp4-small  "
      onClick={() => dispatch({type: 'clearFilter'})}>
      {t('db.clear')}
    </button>
  );
}

function TagFilter() {
  const [tagIsOpen, setTagIsOpen] = useState(false);
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);
  const sortedTagnames = Object.keys(tagData)
    .filter((key) => {
      return key !== '0' && key !== '1' && key !== '2';
    })
    .map((key) => {
      return {
        key: key,
        name: tagData[key]['display_name'],
      };
    });

  return (
    <div className="w-full  overflow-x-hidden no-scrollbar">
      <button
        className=" bp4-button bp4-intent-primary w-full flex-row flex justify-between items-center "
        onClick={() => setTagIsOpen(!tagIsOpen)}>
        <div className=" grow">{t('db.tags')}</div>

        <div className="">{tagIsOpen ? '-' : '+'}</div>
      </button>
      <Collapse isOpen={tagIsOpen}>
        <div className="grid grid-cols-3 gap-2 mt-2 bg-gray-800 p-1">
          {sortedTagnames.map((t) => (
            <TagFilterButton key={t.key} name={t.name} tag={t.key} />
          ))}
        </div>
      </Collapse>
    </div>
  );
}

function TagFilterButton({tag, name}: {tag; name: string}) {
  const filter = useContext(FilterContext);
  const dispatch = useContext(FilterDispatchContext);

  const handleClick = () => {
    dispatch({
      type: 'handleTag',
      tag: tag,
    });
  };
  let intent;
  switch (filter.tagFilter[tag].state) {
    case ItemFilterState.include:
      intent = Intent.SUCCESS;
      break;
    case ItemFilterState.exclude:
      intent = Intent.DANGER;
      break;
    default:
      intent = Intent.NONE;
      break;
  }

  return (
    <Button intent={intent} onClick={handleClick}>
      <div className="text-center">{name}</div>
    </Button>
  );
}

function CharacterFilter() {
  const [charIsOpen, setCharIsOpen] = useState(false);
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);
  const sortedCharNames = charNames.sort((a, b) => {
    if (t(a) < t(b)) {
      return -1;
    }
    if (t(a) > t(b)) {
      return 1;
    }
    return 0;
  });
  const [charSearch, setCharSearch] = useState<string>('');

  const translateCharName = (charName: string) =>
    t('game:character_names.' + charName);

  return (
    <div className="w-full  overflow-x-hidden no-scrollbar">
      <button
        className=" bp4-button bp4-intent-primary w-full flex-row flex justify-between items-center "
        onClick={() => setCharIsOpen(!charIsOpen)}>
        <div className=" grow">{t('db.characters')}</div>

        <div className="">{charIsOpen ? '-' : '+'}</div>
      </button>
      <Collapse isOpen={charIsOpen}>
        <div className="flex flex-col mt-2 bg-gray-800 p-1">
          <label
            htmlFor="email"
            className="relative text-gray-400 focus-within:text-gray-600 flex flex-row">
            <FaSearch className="pointer-events-none w-4 h-4 absolute top-2 transform   right-2 " />

            <input
              className="bp4-input bp4-icon bp4-icon-filter grow"
              type="text"
              dir="auto"
              onChange={(e) => {
                setCharSearch(e.target.value);
              }}
            />
          </label>

          <div className="grid grid-cols-4 gap-1 mt-1 overflow-y-auto overflow-x-hidden">
            {sortedCharNames
              .filter((charName) =>
                translateCharName(charName)
                  .toLocaleLowerCase()
                  .includes(charSearch.toLocaleLowerCase()),
              )
              .map((charName) => (
                <CharFilterButton key={charName} charName={charName} />
              ))}
          </div>
        </div>
      </Collapse>
    </div>
  );
}

function CharFilterButton({charName}: {charName: string}) {
  const filter = useContext(FilterContext);
  const dispatch = useContext(FilterDispatchContext);

  const handleClick = () => {
    dispatch({
      type: 'handleChar',
      char: charName,
    });
  };

  switch (filter.charFilter[charName].state) {
    case ItemFilterState.include:
      return (
        <button
          className={'bp4-button bp4-intent-success block'}
          onClick={handleClick}>
          <CharFilterButtonChild charName={charName} />
        </button>
      );
    case ItemFilterState.exclude:
      return (
        <button
          className={'bp4-button bp4-intent-danger block'}
          onClick={handleClick}>
          <CharFilterButtonChild charName={charName} />
        </button>
      );
    case ItemFilterState.none:
    default:
      return (
        <button className={'bp4-button block '} onClick={handleClick}>
          <CharFilterButtonChild charName={charName} />
        </button>
      );
  }
}

function CharFilterButtonChild({charName}: {charName: string}) {
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);
  const displayCharName = t('game:character_names.' + charName);

  const travelerName = (
    charName.includes('lumine') || charName.includes('aether')
      ? displayCharName
      : ''
  ).replace(/.*?\((\S+)\).*?/, '$1');

  return (
    <div className="flex flex-col truncate gap-1">
      <img
        alt={displayCharName}
        src={`/api/assets/avatar/${charName}.png`}
        className="truncate h-16 object-contain"
      />
      {travelerName != '' ? (
        <div className="text-center">{travelerName}</div>
      ) : (
        <></>
      )}
    </div>
  );
}

function SortBy() {
  const [sortIsOpen, setSortIsOpen] = useState(false);
  const {t: translation} = useTranslation();
  const t = (s: string) => translation<string>(s);

  return (
    <div className="w-full  overflow-x-hidden no-scrollbar">
      <button
        className=" bp4-button bp4-intent-primary w-full flex-row flex justify-between items-center "
        onClick={() => setSortIsOpen(!sortIsOpen)}>
        <div className=" grow">{t('db.sort_by')}</div>

        <div className="">{sortIsOpen ? '-' : '+'}</div>
      </button>
      <Collapse isOpen={sortIsOpen}>
        <div className="flex flex-col mt-2 bg-gray-800 p-1">
          <div className="flex flex-row gap-4">
            {sortByParams.map((param) => (
              <SortByParamButton
                key={param.sortKey}
                sortKey={param.sortKey}
                translation={t(param.translationKey)}
              />
            ))}
          </div>
        </div>
      </Collapse>
    </div>
  );
}

function SortByParamButton({
  sortKey,
  translation,
}: {
  sortKey: string;
  translation: string;
}) {
  const filter = useContext(FilterContext);
  const dispatch = useContext(FilterDispatchContext);

  const handleClick = () => {
    dispatch({
      type: 'handleSortBy',
      sortByKey: sortKey,
    });
  };

  let intent: Intent;
  if (filter.sortBy.sortKey !== sortKey) {
    intent = 'none';
  } else {
    switch (filter.sortBy.sortByDirection) {
      case SortByDirection.asc:
        intent = 'success';
        break;
      case SortByDirection.dsc:
        intent = 'danger';
        break;
      default:
        intent = 'none';
        break;
    }
  }

  return (
    <Button onClick={handleClick} intent={intent}>
      <div className="flex flex-row gap-1 justify-center items-center">
        {intent === 'success' && <FaArrowUp />}
        {intent === 'danger' && <FaArrowDown />}
        {translation}
      </div>
    </Button>
  );
}
