import { Callout, useHotkeys } from '@blueprintjs/core';
import React from 'react';

import { SectionDivider } from '~src/Components/SectionDivider';
import { Viewport } from '~src/Components/Viewport';
import { RootState, useAppDispatch, useAppSelector } from '~src/store';
import { simActions } from '..';
import { ActionList, SimProgress } from '../Components';
import { Team } from './Team';
import { Trans, useTranslation } from 'react-i18next';
import { Toolbox } from './Toolbox';
import { ActionListTooltip, TeamBuilderTooltip } from './Tooltips';
import { updateCfg, ready } from '../simSlice';

export function Simple() {
  let { t } = useTranslation();

  const { cfg, cfg_err, showBuilder } = useAppSelector((state: RootState) => {
    return {
      cfg: state.sim.cfg,
      cfg_err: state.sim.cfg_err,
      showBuilder: state.sim.showBuilder,
    };
  });
  const dispatch = useAppDispatch();

  const [open, setOpen] = React.useState<boolean>(false);

  const hotkeys = React.useMemo(
    () => [
      {
        combo: 'Esc',
        global: true,
        label: t('simple.exit_edit'),
        onKeyDown: () => {
          dispatch(simActions.editCharacter({ index: -1 }));
        },
      },
    ],
    []
  );
  useHotkeys(hotkeys);

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col gap-2">
        <div className="flex flex-col">
          {showBuilder ? (
            <>
              <SectionDivider>
                <Trans>simple.team</Trans>
              </SectionDivider>
              <TeamBuilderTooltip />
              <Team />
            </>
          ) : null}

          <SectionDivider>
            <Trans>simple.action_list</Trans>
          </SectionDivider>

          <ActionListTooltip />

          <ActionList cfg={cfg} onChange={(v) => dispatch(updateCfg(v))} />

          <div className="sticky bottom-0 bg-bp-bg flex flex-col gap-y-1">
            {cfg_err !== '' ? (
              <div className="pl-2 pr-2 pt-2 mt-1">
                <Callout intent="warning" title="Error parsing config">
                  <pre className=" whitespace-pre-wrap">{cfg_err}</pre>
                </Callout>
              </div>
            ) : null}
            <div className="pl-2 pr-2 pt-2 mt-1">
              <Callout intent="warning" title="Breaking changes">
                Please be aware that there have been syntax changes with the
                core rewrite. Your existing configs may not work. Please check
                out the{' '}
                <a href="https://docs.gcsim.app/migration" target="_blank">
                  migration guide
                </a>
              </Callout>
            </div>
            <Toolbox canRun={cfg_err === '' && ready()} />
          </div>
        </div>
      </div>
      <SimProgress isOpen={open} onClose={() => setOpen(false)} />
    </Viewport>
  );
}
