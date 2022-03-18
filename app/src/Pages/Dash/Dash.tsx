import { Card, Elevation, Icon } from "@blueprintjs/core";
import { useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";

interface DashCardProps {
  children: React.ReactNode;
  onClick?: () => void;
}

function DashCard({ children, onClick }: DashCardProps) {
  return (
    <div className="main-page-button-container">
      <Card
        interactive
        elevation={Elevation.TWO}
        className="main-page-card"
        onClick={onClick}
      >
        {children}
      </Card>
    </div>
  );
}

export function Dash() {
  useTranslation()

  const [_, setLocation] = useLocation();
  return (
    <main className="w-full flex flex-col items-center flex-grow ">
      <span className="font-bold text-md mt-4 p-1">
        <a href="https://github.com/genshinsim/gcsim" target="_blank">
          gcsim
        </a>{" "}
        <Trans>dash.is_a_team</Trans>
      </span>
      <div className="flex flex-row flex-initial flex-wrap w-full lg:w-[60rem] mt-4">
        <DashCard onClick={() => setLocation("/simple")}>
          <span className="font-bold text-xl">
            <Icon icon="calculator" className="mr-2" size={25} />
            <Trans>dash.simulator</Trans>
          </span>
        </DashCard>

        <DashCard onClick={() => setLocation("/advanced")}>
          <span className="font-bold text-xl">
            <Icon icon="rocket-slant" className="mr-2" size={25} />
            <Trans>dash.advanced_mode</Trans>
          </span>
        </DashCard>

        <DashCard onClick={() => setLocation("/viewer")}>
          <span className="font-bold text-xl">
            <Icon icon="chart" className="mr-2" size={25} />
            <Trans>dash.viewer</Trans>
          </span>
        </DashCard>

        <DashCard onClick={() => setLocation("/db")}>
          <span className="font-bold text-xl">
            <Icon icon="database" className="mr-2" size={25} />
            <Trans>dash.teams_db</Trans>
          </span>
        </DashCard>

        <DashCard
          onClick={() =>
            window.open(
              "https://github.com/genshinsim/gcsim/releases",
              "_blank"
            )
          }
        >
          <span className="font-bold text-xl">
            <Icon icon="download" className="mr-2" size={25} />
            <Trans>dash.desktop_tool</Trans>
          </span>
        </DashCard>

        <DashCard
          onClick={() => window.open("https://docs.gcsim.app", "_blank")}
        >
          <span className="font-bold text-xl">
            <Icon icon="document" className="mr-2" size={25} />
            <Trans>dash.documentation</Trans>
          </span>
        </DashCard>

        <DashCard
          onClick={() =>
            window.open(
              "https://github.com/genshinsim/gcsim#Contributing",
              "_blank"
            )
          }
        >
          <span className="font-bold text-xl">
            <Icon icon="git-branch" className="mr-2" size={25} />
            <Trans>dash.contribute</Trans>
          </span>
        </DashCard>

        <DashCard onClick={() => setLocation("/about")}>
          <span className="font-bold text-xl">
            <Icon icon="info-sign" className="mr-2" size={25} />
            <Trans>dash.about</Trans>
          </span>
        </DashCard>
      </div>
    </main>
  );
}
