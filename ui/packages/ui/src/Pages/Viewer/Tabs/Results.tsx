import Overview from "../Components/Results/Overview";
import { SimResults } from "@gcsim/types";
import TeamHeader from "../Components/Results/TeamHeader";
import DistributionCard from "../Components/Results/DistributionCard";
import { Card } from "@blueprintjs/core";
import { ReactNode } from "react";
import classNames from "classnames";
import CharacterDPSCard from "../Components/Damage/CharacterDPSCard";
import ElementDPSCard from "../Components/Damage/ElementDPSCard";
import TargetDPSCard from "../Components/Damage/TargetDPSCard";

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  return (
    <div className="w-full 2xl:mx-auto 2xl:container px-2">
      {/* Overview */}
      <Group>
        <TeamHeader data={data} />
        <Overview data={data} />
        <Card className="flex col-span-3 h-24 min-h-full">
          Target info + sim metadata (num iterations)
        </Card>
        <DistributionCard data={data} />
      </Group>

      {/* Damage */}
      <Group>
        <Heading text="Damage" />
        <Card className="col-span-full h-64">
          Damage Timeline (dps/time + cumu %)
        </Card>
        <CharacterDPSCard data={data} />
        <ElementDPSCard data={data} />
        <TargetDPSCard data={data} />
        <Card className="flex col-span-full h-64 min-h-full">
          DPS per character w/ (element breakdown & target breakdown)
        </Card>
        <Card className="flex col-span-full h-64 min-h-full">
          Damage breakdown table(s)
        </Card>
      </Group>

      {/* Energy */}
      <Group>
        <Heading text="Energy" />
        <Card className="flex col-span-full h-64 min-h-full">
          Energy over time + cumu gained + cumu wasted
        </Card>
        <Card className="flex col-span-full h-64 min-h-full">
          Energy produced by source
        </Card>
        <Card className="flex col-span-full h-64 min-h-full">
          Incoming energy per character breakdown
        </Card>
      </Group>

      {/* Reactions & Auras */}
      <Group>
        <Heading text="Reactions & Auras" />
        <Card className="flex col-span-4 h-64 min-h-full">
          Aura uptime timeline (worst, best, heatmap)
        </Card>
        <Card className="flex col-span-2 h-64 min-h-full">
          Aura uptime (pie? vs bar?)
        </Card>
        <Card className="flex col-span-3 h-64 min-h-full">
          Reactions triggered bar chart
        </Card>
        <Card className="flex col-span-3 h-64 min-h-full">
          Reactions by source?
        </Card>
      </Group>

      {/* Healing */}
      {/* TODO: make optional? */}
      <Group>
        <Heading text="Healing" />
        <Card className="flex col-span-full h-64 min-h-full">
          Effective Healing Timeline (+ HP per char?)
        </Card>
        <Card className="flex col-span-3 h-64 min-h-full">
          Healing by Src
        </Card>
        <Card className="flex col-span-3 h-64 min-h-full">
          Healing by target
        </Card>
      </Group>

      {/* Shields */}
      {/* TODO: make optional? */}
      <Group>
        <Heading text="Shields" />
        <Card className="flex col-span-full h-64 min-h-full">
          Shield timeline (worst, best, heatmap)
        </Card>
        <Card className="flex col-span-3 h-[512px] min-h-full">
          Shield uptime bar chart + pie
        </Card>
        <Card className="flex col-span-3 h-[512px] min-h-full">
          Shield hp bar chart
        </Card>
        <Card className="flex col-span-full h-64 min-h-full">
          Shield table
        </Card>
      </Group>

      {/* Sim Metadata? */}
      <Group>
        <Heading text="Simulation Details" />
        <Card className="flex col-span-2 h-64 min-h-full">
          Character uptime (pie?)
        </Card>
        <Card className="flex col-span-4 h-64 min-h-full">
          Character Uptime timeline (worst, best, heatmap)
        </Card>
        <Card className="flex col-span-4 h-64 min-h-full">
          Failed actions bar graph
        </Card>
        <Card className="flex col-span-2 h-64 min-h-full">
          Faied actions timeline (worst, best, heatmap)
        </Card>
        {/* tables? */}
      </Group>
    </div>
  );
};

const Heading = ({ text }: { text: string }) => (
  <h2 className="group flex whitespace-pre-wrap col-span-full text-xl font-semibold mt-8 mb-2">
    {text}
    {/* currently does not work (wouter doesn't support hash links) */}
    {/* <a
        href={target}
        className="ml-2 text-blue-500 opacity-0 transition-opacity group-hover:opacity-100">
      #
    </a> */}
  </h2>
);

type GroupProps = {
  children: ReactNode;
  className?: string;
}

const Group = ({ children, className }: GroupProps) => {
  const cls = classNames(
      className,
      "grid overflow-hidden",
      "grid-cols-1 sm:grid-cols-6",
      "gap-y-2", "sm:gap-2");

  return (
    <div className={cls}>
      {children}
    </div>
  );
};