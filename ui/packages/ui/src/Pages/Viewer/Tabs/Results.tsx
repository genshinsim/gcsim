import { SimResults } from "@gcsim/types";
import { Card, Colors } from "@blueprintjs/core";
import { ReactNode, useEffect, useRef } from "react";
import classNames from "classnames";
import { DistributionCard, RollupCards, TargetInfo, TeamHeader } from "../Components/Overview";
import { CharacterDPSBarChart, CharacterDPSCard, DamageTimelineCard, ElementDPSCard, TargetDPSCard } from "../Components/Damage";
import { useLocation } from "react-router";
import { FiLink2 } from "react-icons/fi";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export default (props: Props) => {
  useScrollToLocation();

  return (
    <div className="w-full 2xl:mx-auto 2xl:container px-2">
      <Overview {...props} />
      <Damage {...props} />
      {/* <Energy {...props} />
      <Reactions {...props} />
      <Healing {...props} />
      <Shields {...props} />
      <SimDetails {...props} /> */}
    </div>
  );
};

const Overview = ({ data }: Props) => (
  <Group>
    <TeamHeader characters={data?.character_details} />
    <RollupCards data={data} />
    <TargetInfo data={data} />
    <DistributionCard data={data} />
  </Group>
);

const Damage = ({ data, running, names }: Props) => (
  <Group>
    <Heading text="Damage" target="damage" color={Colors.VERMILION5} />
    <DamageTimelineCard data={data} running={running} names={names} />

    <CharacterDPSCard data={data} running={running} names={names} />
    <ElementDPSCard data={data} running={running} />
    <TargetDPSCard data={data} running={running} />

    <CharacterDPSBarChart data={data} running={running} names={names} />

    {/* <Card className="flex col-span-full h-64 min-h-full">
      Damage breakdown table(s)
    </Card> */}
  </Group>
);

const Energy = ({ }: Props) => (
  <Group>
    <Heading text="Energy" target="energy" color={Colors.CERULEAN5} />
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
);

const Reactions = ({ }: Props) => (
  <Group>
    <Heading text="Reactions & Auras" target="reactions" color={Colors.VIOLET5} />
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
);

const Healing = ({ }: Props) => (
  <Group>
    <Heading text="Healing" target="healing" color={Colors.FOREST5} />
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
);

const Shields = ({ }: Props) => (
  <Group>
    <Heading text="Shields" target="shields" color={Colors.GOLD5} />
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
);

const SimDetails = ({ }: Props) => (
  <Group>
    <Heading text="Simulation Details" target="sim" color={Colors.TURQUOISE5} />
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
);

type HeadingProps = {
  text: string;
  target: string;
  color?: string;
}

const Heading = ({ text, target, color }: HeadingProps) => {
  const linkClass = classNames(
    "ml-3 mt-1",
    "text-blue-500",
    "opacity-0 group-hover:opacity-100 transition-opacity",
    "flex justify-center items-center"
  );

  return (
    <h2 className="group flex whitespace-pre-wrap col-span-full text-xl font-semibold mt-12 mb-1">
      <span style={{ color: color }} >{text}</span>
      <a href={"#" + target} id={target} className={linkClass}>
        <FiLink2 />
      </a>
    </h2>
  );
};

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

function useScrollToLocation() {
  const scrolled = useRef(false);
  const { key, hash } = useLocation();
  const prevKey = useRef(key);

  useEffect(() => {
    if (hash == null) {
      return;
    }

    if (prevKey.current !== key) {
      prevKey.current = key;
      scrolled.current = false;
    }

    if (!scrolled.current) {
      const id = hash.replace('#', '');
      const element = document.getElementById(id);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
        scrolled.current = true;
      }
    }
  });
}