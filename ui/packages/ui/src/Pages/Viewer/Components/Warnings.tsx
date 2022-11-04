import DismissibleCallout from "../../../Components/DismissibleCallout";
import { Intent } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { i18n } from "i18next";

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

function fmt(i18n: i18n, val?: number) {
  return val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) + "s";
}

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
  const { i18n } = useTranslation();

  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_energy ?? false);

  const characters = data?.statistics?.failed_actions?.map((fa, i) => {
    if (data.character_details == null) {
      return null;
    }

    return (
      <>
        <div key={i} className="list-item">{data.character_details[i].name}</div>
        <div>{fmt(i18n, fa.insufficient_energy)}</div>
      </>
    );
  });

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
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">avg insufficient energy duration</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-2 justify-start">
          {characters}
        </ul>
      </div>
    </DismissibleCallout>
  );
};

const SwapWarning = ({ data }: WarningProps) => {
  const { i18n } = useTranslation();

  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.swap_cd ?? false);

  const characters = data?.statistics?.failed_actions?.map((fa, i) => {
    if (data.character_details == null) {
      return null;
    }

    return (
      <>
        <div key={i} className="list-item">{data.character_details[i].name}</div>
        <div>{fmt(i18n, fa.swap_cd)}</div>
      </>
    );
  });

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
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">avg swap cd duration</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-2 justify-start">
          {characters}
        </ul>
      </div>
    </DismissibleCallout>
  );
};

const CooldownWarning = ({ data }: WarningProps) => {
  const { i18n } = useTranslation();

  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.skill_cd ?? false);

  const characters = data?.statistics?.failed_actions?.map((fa, i) => {
    if (data.character_details == null) {
      return null;
    }

    return (
      <>
        <div key={i} className="list-item">{data.character_details[i].name}</div>
        <div>{fmt(i18n, fa.skill_cd)}</div>
      </>
    );
  });

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
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">avg skill cd duration</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-2 justify-start">
          {characters}
        </ul>
      </div>
    </DismissibleCallout>
  );
};

const StaminaWarning = ({ data }: WarningProps) => {
  const { i18n } = useTranslation();

  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_stamina ?? true);

  const characters = data?.statistics?.failed_actions?.map((fa, i) => {
    if (data.character_details == null) {
      return null;
    }

    return (
      <>
        <div key={i} className="list-item">{data.character_details[i].name}</div>
        <div>{fmt(i18n, fa.insufficient_stamina)}</div>
      </>
    );
  });

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
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">avg insufficient stamina duration</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-2 justify-start">
          {characters}
        </ul>
      </div>
    </DismissibleCallout>
  );
};