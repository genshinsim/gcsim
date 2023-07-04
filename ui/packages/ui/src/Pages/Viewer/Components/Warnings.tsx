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
        An abnormal amount of iterations failed to execute a burst because the character did not
        have enough energy. Consider updating the config to better manage energy, as no actions are
        performed during failures.
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
        An abnormal amount of iterations failed to execute a swap because swap was on cooldown.
        Consider updating the config to better account for swap cooldowns, as no actions are
        performed during failures.
      </p>
      <FailedActionDetails
          title="swap cd failure duration"
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
        An abnormal amount of iterations failed to execute a skill because the skill was on cooldown.
        Consider updating the config to better account for skill cooldowns, as no actions are
        performed during failures.
      </p>
      <FailedActionDetails
          title="skill cd failure duration"
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
        An abnormal amount of iterations failed to execute an action because the character did not
        have enough stamina. Consider updating the config to better manage stamina, as no actions
        are performed during failures.
      </p>
      <FailedActionDetails
          title="insufficient stamina duration"
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
        title="Unable to Use Character Dash (Dash on CD)"
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        An abnormal amount of iterations failed to execute a dash because the dash was on cooldown.
        Consider updating the config to better account for dash cooldowns, as no actions are
        performed during failures.
      </p>
      <FailedActionDetails
          title="dash cd failure duration"
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