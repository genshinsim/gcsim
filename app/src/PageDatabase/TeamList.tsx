import { AnchorButton, Button, ButtonGroup, H3 } from '@blueprintjs/core';
import React from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';
import { Link, useLocation } from 'wouter';
import { TeamCard } from '~src/Components/Team';
import { updateCfg } from '~src/Pages/Sim';
import { useAppSelector } from '~src/store';
import { loadCharacter } from './dbSlice';
import { SendConfigDialog } from './SendConfigDialog';

type TeamListProps = {
  char: string;
  team: string;
};

export function TeamsList({ char, team }: TeamListProps) {
  useTranslation();
  const [config, setConfig] = React.useState<string>('');
  const charSims = useAppSelector((state) => state.db.charSims);
  const dispatch = useDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    if (!(char in charSims)) {
      dispatch(loadCharacter(char));
    }
  }, [charSims, dispatch]);

  if (!(char in charSims) || charSims[char].length === 0) {
    return (
      <div className="flex flex-row place-content-center mt-2">
        Sorry, <span className=" capitalize">{team.replaceAll('-', ', ')}</span>
        does not have any simulations yet
      </div>
    );
  }

  const handleOpenInSim = (keep: boolean) => {
    dispatch(updateCfg(config, keep));
    setLocation('/simulator');
    setConfig('');
  };

  const sims = charSims[char].filter((s) => {
    const key = s.metadata.char_names
      .map((x) => x)
      .sort()
      .join('-');
    return team === key;
  });

  const rows = sims.map((s) => {
    const details = (
      <>
        <span>
          <Trans>db.total_dps</Trans>
          {parseInt(s.metadata.dps.mean.toFixed(0)).toLocaleString()}
        </span>
        <span>
          <Trans>db.number_of_targets</Trans>
          {s.metadata.num_targets}
        </span>
        <span>
          <Trans>db.average_dps_per</Trans>
          {parseInt(
            (s.metadata.dps.mean / s.metadata.num_targets).toFixed(0)
          ).toLocaleString()}
        </span>
      </>
    );

    const action = (
      <ButtonGroup vertical>
        <Link href={'/v3/viewer/share/' + s.simulation_key}>
          <AnchorButton small rightIcon="chart">
            <Trans>db.show_in_viewer</Trans>
          </AnchorButton>
        </Link>
        <Button
          small
          rightIcon="rocket-slant"
          onClick={() => setConfig(s.config)}
        >
          <Trans>db.load_in_simulator</Trans>
        </Button>
      </ButtonGroup>
    );

    return (
      <TeamCard
        meta={s.metadata}
        key={s.simulation_key}
        summary={details}
        actions={action}
        onCharacterClick={(char) => setLocation(`/db/${char}`)}
      />
    );
  });

  return (
    <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
      <H3>
        Showing Teams for{' '}
        <span className=" capitalize">{team.replaceAll('-', ', ')}</span>
      </H3>
      <div className="flex flex-col">{rows}</div>
      <SendConfigDialog
        isOpen={config !== ''}
        handleClose={() => setConfig('')}
        handleSubmit={handleOpenInSim}
      />
    </main>
  );
}
