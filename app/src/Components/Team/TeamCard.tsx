import React from 'react';
import { Popover2, Tooltip2 } from '@blueprintjs/popover2';
import { Metadata } from '~src/Types/stats';
import { CharacterTooltip } from '~src/Components/Character';
import { Trans, useTranslation } from 'react-i18next';

type TeamCardProps = {
  meta: Metadata;
  summary?: React.ReactNode;
  actions?: React.ReactNode;
};

export function TeamCard({
  meta,
  summary = null,
  actions = null,
}: TeamCardProps) {
  useTranslation();
  const chars = meta.char_details.map((char) => {
    return (
      <Popover2>
        <Tooltip2 content={<CharacterTooltip char={char} />}>
          <div className="hover:bg-gray-600 border border-gray-700 hover:border-gray-400 rounded-md relative">
            <img
              src={'/images/avatar/' + char.name + '.png'}
              alt={char.name}
              className="w-16"
              key={char.name}
            />
            <div className=" absolute top-0 right-0 text-sm font-semibold text-grey-300">{`${char.cons}`}</div>
          </div>
        </Tooltip2>
      </Popover2>
    );
  });

  return (
    <div className="flex flex-row flex-wrap sm:flex-nowrap gap-y-1 w-full m-2 p-2 rounded-md bg-gray-700 place-items-center">
      <div className="flex flex-col sm:basis-1/4 xs:basis-full">
        <div className="grid grid-cols-4">{chars}</div>
        <div className="hidden basis-0 lg:block md:flex-1"></div>
      </div>
      <div className=" flex-1 overflow-hidden mb-auto pl-2 hidden lg:block"></div>
      <div className="ml-auto flex flex-col mr-4 md:basis-60 basis-full">
        {summary}
      </div>
      <div>{actions}</div>
    </div>
  );
}
