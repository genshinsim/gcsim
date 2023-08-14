import DismissibleCallout from "../../../Components/DismissibleCallout";
import { Intent } from "@blueprintjs/core";
import { FailedActions, FloatStat, SimResults } from "@gcsim/types";
import { useState } from "react";
import { useTranslation } from "react-i18next";

type WarningProps = {
  data: SimResults | null;
}

// TODO: translation
export default (props: WarningProps) => {
  const warnings = [
    <IncompleteCharWarning key="incomplete" {...props} />,
    <PositionOverlapWarning key="target" {...props} />,
    <EnergyWarning key="energy" {...props} />,
    <CooldownWarning key="cd" {...props} />,
    <StaminaWarning key="stamina" {...props} />,
    <SwapWarning key="swap" {...props} />,
    <DashWarning key="dash" {...props} />
  ];

  return (
    <div className="flex flex-col gap-2 pt-4 empty:pt-0 w-full px-2 2xl:mx-auto 2xl:container">
      {warnings}
    </div>
  );
};

const IncompleteCharWarning = ({ data }: WarningProps) => {  
  const [show, setShow] = useState(true);
  const incomplete = data?.incomplete_characters;
  const visible = show && (incomplete != null && incomplete.length > 0);
  
  const link = (
    <a href="https://discord.gg/m7jvjdxx7q" target="_blank" rel="noreferrer">
      gcsim discord!
    </a>
  );

  return (
    <DismissibleCallout
        title="Incomplete Characters Used"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        This simulation contains early release characters! These characters are fully implemented,
        but may not have optimal frame data aligned with in-game animations. We are actively
        collecting data to improve their implementation. If you wish to help,
        please reach out in the {link}
      </p>
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">incomplete characters</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-3 justify-start">
          {data?.incomplete_characters?.map(c => (
            <div key={c} className="list-item">{c}</div>
          ))}
        </ul>
      </div>
    </DismissibleCallout>
  );
};

const PositionOverlapWarning = ({ data }: WarningProps) => {  
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.target_overlap ?? false);
  
  return (
    <DismissibleCallout
        title="Target Positions Overlap"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {"Target position's overlap in at least on iteration. Confirm if this is intended and update positions to avoid overlaps as necessary. Overlapping positions may result in inaccurate simulations."}
      </p>
    </DismissibleCallout>
  );
};

const EnergyWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_energy ?? false);

  return (
    <DismissibleCallout
        title="Delay in Burst - Potential Energy Deficiency"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Some iterations delayed executing one or more bursts due to lack of energy. This causes the active character to idle until enough energy is gained (see <a href="https://docs.gcsim.app/guides/understanding_config_files#gcsim-script-gcsl">here</a>). Consider updating the config if the downtime below is undesired. 
      </p>
      <FailedActionDetails
          title="total burst delay duration per iteration"
          data={data}
          stat={(fa) => fa.insufficient_energy} />
    </DismissibleCallout>
  );
};

const SwapWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.swap_cd ?? false);

  return (
    <DismissibleCallout
        title="Delay in Swapping - Swap on CD"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Some iterations delayed executing one or more swaps due to its cooldown. This causes the active character to idle until swap is off cooldown. Consider updating the config if the downtime below is undesired.
      </p>
      <FailedActionDetails
          title="total swap delay due to cd per iteration"
          data={data}
          stat={(fa) => fa.swap_cd} />
    </DismissibleCallout>
  );
};

const CooldownWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.skill_cd ?? false);

  return (
    <DismissibleCallout
        title="Delay in Skill - Skill on CD"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Some iterations delayed executing one or more skills due to its cooldown. This causes the active character to idle until their skill is off cooldown. Consider updating the config if the downtime below is undesired.
      </p>
      <FailedActionDetails
          title="total skill delay due to cd per iteration"
          data={data}
          stat={(fa) => fa.skill_cd} />
    </DismissibleCallout>
  );
};

const StaminaWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_stamina ?? false);

  return (
    <DismissibleCallout
        title="Delay in Dash - Insufficient Stamina"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Some iterations delayed executing dash due to insufficient stamina. This causes the active character to idle until enough stamina regnerated. Consider updating the config if the downtime below is undesired.
      </p>
      <FailedActionDetails
          title="total delay due to insufficient stamina per iteration"
          data={data}
          stat={(fa) => fa.insufficient_stamina} />
    </DismissibleCallout>
  );
};

const DashWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.dash_cd ?? false);

  return (
    <DismissibleCallout
        title="Delay in Dash - Dash on CD"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Some iterations delayed executing dash due to its cooldown. This causes the active character to idle until enough stamina regnerated. Consider updating the config if the downtime below is undesired.
      </p>
      <FailedActionDetails
          title="total dash delay due to cd per iteration"
          data={data}
          stat={(fa) => fa.dash_cd} />
    </DismissibleCallout>
  );
};

type DetailsProps = {
  data: SimResults | null;
  title: string;
  stat: (x: FailedActions) => FloatStat | undefined;
}

const FailedActionDetails = ({ data, title, stat }: DetailsProps) => {
  const { i18n } = useTranslation();

  if (data?.character_details == null) {
    return null;
  }

  function fmt(val?: number) {
    return val?.toLocaleString(
        i18n.language, { maximumFractionDigits: 2 }) + "s";
  }

  const Item = ({ f, i }: { f: FloatStat | undefined, i: number }) => {
    if (f?.max == 0) {
      return null;
    }

    return (
      <>
        <div className="list-item">{data.character_details?.[i].name}</div>
        <div>
          {
            "avg: " + fmt(f?.mean)
            + " | min: " + fmt(f?.min)
            + " | max: " + fmt(f?.max)
            + " | std: " + fmt(f?.sd)
          }
        </div>
      </>
    );
  };

  const details = data?.statistics?.failed_actions?.map((fa, i) => (
    <Item key={i.toString()} f={stat(fa)} i={i} />
  ));

  return (
    <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
      <span className="font-bold">{title}</span>
      <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-3 justify-start">
        {details}
      </ul>
    </div>
  );
};