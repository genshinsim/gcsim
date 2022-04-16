import { Card, Elevation, Icon } from "@blueprintjs/core";
import { Link } from "wouter";
import { Trans, useTranslation } from "react-i18next";

interface DashCardProps {
  children: React.ReactNode;
  href: string;
  target?: string;
}

function DashCard({ children, href, target }: DashCardProps) {
  return (
    <div className="main-page-button-container">
      {target?
        <a href={href} target={target}>
          <Card
            interactive
            elevation={Elevation.TWO}
            className="main-page-card"
          >
            {children}
          </Card>
        </a>
        :
        <Link href={href}><a>
          <Card
            interactive
            elevation={Elevation.TWO}
            className="main-page-card"
          >
            {children}
          </Card>
        </a></Link>
      }
    </div>
  );
}

export function Dash() {
  useTranslation();
  return (
    <main className="w-full flex flex-col items-center flex-grow ">
      <span className="font-bold text-md mt-4 p-1">
        <a href="https://github.com/genshinsim/gcsim" target="_blank">
          gcsim
        </a>{" "}
        <Trans>dash.is_a_team</Trans>
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

        <DashCard href="/db">
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

        <DashCard
          href="https://docs.gcsim.app"
          target="_blank"
        >
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
    </main >
  );
}
