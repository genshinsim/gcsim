import React from 'react';
import { Callout, Card, Elevation, Icon } from '@blueprintjs/core';
import { Trans, useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

interface DashCardProps {
  children: React.ReactNode;
  href: string;
  target?: string;
}

function DashCard({ children, href, target }: DashCardProps) {
  return (
    <div className="main-page-button-container">
      {target ? (
        <a href={href} target={target}>
          <Card
            interactive
            elevation={Elevation.TWO}
            className="main-page-card"
          >
            {children}
          </Card>
        </a>
      ) : (
        <Link to={href}>
          <a>
            <Card
              interactive
              elevation={Elevation.TWO}
              className="main-page-card"
            >
              {children}
            </Card>
          </a>
        </Link>
      )}
    </div>
  );
}

export function Dash() {
  useTranslation();
  return (
    <main className="w-full flex flex-col items-center flex-grow pb-4">
      <span>
        <Callout intent="success" className=" max-w-[600px] mt-4">
          Thank you for your patience. The core rewrite is now complete. Hitlag
          has been implemented along with a ton of config syntax improvements.
          <br />
          <div className="mt-2 font-bold">
            Please check out the migration guide here:{' '}
            <a href="https://docs.gcsim.app/migration" target="_blank" rel="noreferrer">
              Migration Guide
            </a>
          </div>
        </Callout>
      </span>
      <div className="flex flex-row flex-initial flex-wrap w-full lg:w-[60rem] mt-4">
        <DashCard href="/simulator">
          <span className="font-bold text-xl">
            <Icon icon="calculator" className="mr-2" size={25} />
            <Trans>dash.simulator</Trans>
          </span>
        </DashCard>

        <DashCard href="/viewer">
          <span className="font-bold text-xl">
            <Icon icon="chart" className="mr-2" size={25} />
            <Trans>dash.viewer</Trans>
          </span>
        </DashCard>

        <DashCard href="https://db.gcsim.app" target="_blank">
          <span className="font-bold text-xl">
            <Icon icon="database" className="mr-2" size={25} />
            <Trans>dash.teams_db</Trans>
          </span>
        </DashCard>

        <DashCard
          href="https://github.com/genshinsim/gcsim/releases"
          target="_blank"
        >
          <span className="font-bold text-xl">
            <Icon icon="download" className="mr-2" size={25} />
            <Trans>dash.desktop_tool</Trans>
          </span>
        </DashCard>

        <DashCard href="https://docs.gcsim.app" target="_blank">
          <span className="font-bold text-xl">
            <Icon icon="document" className="mr-2" size={25} />
            <Trans>dash.documentation</Trans>
          </span>
        </DashCard>

        <DashCard
          href="https://github.com/genshinsim/gcsim#Contributing"
          target="_blank"
        >
          <span className="font-bold text-xl">
            <Icon icon="git-branch" className="mr-2" size={25} />
            <Trans>dash.contribute</Trans>
          </span>
        </DashCard>

        <DashCard href="/about">
          <span className="font-bold text-xl">
            <Icon icon="info-sign" className="mr-2" size={25} />
            <Trans>dash.about</Trans>
          </span>
        </DashCard>
      </div>
    </main>
  );
}
