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
    <PositionOverlapWarning key="target" {...props} />,
    <EnergyWarning key="energy" {...props} />,
    <CooldownWarning key="cd" {...props} />,
    <StaminaWarning key="stamina" {...props} />,
    <SwapWarning key="swap" {...props} />,
  ];

  return (
    <div className="flex flex-col gap-2 pt-4 empty:pt-0 mx-auto max-w-2xl w-full">
      {warnings}
    </div>
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
        {"Target position's overlap in at least one iteration. This may result in inaccurate "
          + "simulations. Update positions to avoid any overlaps."}
      </p>
    </DismissibleCallout>
  );
};

const EnergyWarning = ({ data }: WarningProps) => {
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_energy ?? false);

  return (
    <DismissibleCallout
        title="Insufficient Energy"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        At least one character frequently did not have enough energy throughout the simulation to
        burst. This causes the simulation to wait for more energy, but will perform no actions while
        waiting. Increase ER or update the config to reduce energy requirements.
      </p>
      <FailedActionDetails
          title="insufficient energy duration"
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
        title="Unable to Swap Characters (Swap on CD)"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Character swaps were delayed throughout the simulation due to character swap being on
        cooldown. This causes the simulation to wait for the cooldown to end, but will perform no
        actions while waiting. Update the config to better account for the swap cooldown.
      </p>
      <FailedActionDetails
          title="swap cd duration"
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
        title="Unable to Use Character Skills (Skills on CD)"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        Skills were frequently attempted to be used when on cooldown. This causes the simulation
        to wait for the cooldown to end, but will perform no actions while waiting. Update the
        config to better account for skill cooldowns.
      </p>
      <FailedActionDetails
          title="skill cd duration"
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
        title="Insufficient Stamina"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        At least one character frequently did not have enough stamina throughout the simulation.
        This causes the simulation to wait until there is enough stamina to perform an action, but
        will perform no other actions while waiting. Update the config to better manage stamina.
      </p>
      <FailedActionDetails
          title="insufficient stamina duration"
          data={data}
          stat={(fa) => fa.insufficient_stamina} />
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