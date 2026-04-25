import DismissibleCallout from "../../../Components/DismissibleCallout";
import { Intent } from "@blueprintjs/core";
import { FailedActions, FloatStat, SimResults } from "@gcsim/types";
import { useState } from "react";
import { Trans, useTranslation } from "react-i18next";

type WarningProps = {
  data: SimResults | null;
}

// TODO: translation
export default (props: WarningProps) => {
  const warnings = [
    <IncompleteCharWarning key="incomplete" {...props} />,
    <PositionOverlapWarning key="target" {...props} />,
    <EnergyWarning key="energy" {...props} />,
    <BurstWarningCD key="burst_cd" {...props} />,
    <SkillWarning key="cd" {...props} />,
    <StaminaWarning key="stamina" {...props} />,
    <SwapWarning key="swap" {...props} />,
    <DashWarning key="dash" {...props} />,
    <IgnoreBurstEnergyMode key="ignore_burst_energy" {...props} />
  ];

  return (
    <div className="flex flex-col gap-2 pt-4 empty:pt-0 w-full px-2 2xl:mx-auto 2xl:container">
      {warnings}
    </div>
  );
};

const IncompleteCharWarning = ({ data }: WarningProps) => {  
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const incomplete = data?.incomplete_characters;
  const visible = show && (incomplete != null && incomplete.length > 0);

  return (
    <DismissibleCallout
        title={t<string>("warnings.incomplete_char_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        <Trans i18nKey="warnings.incomplete_char_body">
          <a href="https://discord.gg/m7jvjdxx7q" target="_blank" rel="noreferrer"/>
        </Trans>
      </p>
      <div className="flex flex-col justify-start gap-1 text-xs pt-2 font-mono text-gray-400">
        <span className="font-bold">{t<string>("warnings.incomplete_char_data_header")}</span>
        <ul className="list-disc pl-4 grid grid-cols-[auto_minmax(0,_1fr)] gap-x-3 justify-start">
          {data?.incomplete_characters?.map(c => (
            <div key={c} className="list-item">{t<string>("character_names." + c, { ns: "game" })}</div>
          ))}
        </ul>
      </div>
    </DismissibleCallout>
  );
};

const PositionOverlapWarning = ({ data }: WarningProps) => { 
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.target_overlap ?? false);
  
  return (
    <DismissibleCallout
        title={t<string>("warnings.position_overlap_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.position_overlap_body")}
      </p>
    </DismissibleCallout>
  );
};

const EnergyWarning = ({ data }: WarningProps) => {
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_energy ?? false);

  return (
    <DismissibleCallout
        title={t<string>("warnings.energy_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.energy_body")}
      </p>
      <FailedActionDetails
          title={t<string>("warnings.energy_data_header")}
          data={data}
          stat={(fa) => fa.insufficient_energy} />
    </DismissibleCallout>
  );
};

const BurstWarningCD = ({data}: WarningProps) => {
  const {t} = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.burst_cd ?? false);

  return (
    <DismissibleCallout
      title={t<string>("warnings.burst_cd_title")}
      intent={Intent.WARNING}
      show={visible}
      onDismiss={() => setShow(false)}>
      <p>{t<string>("warnings.burst_cd_body")}</p>
      <FailedActionDetails
        title={t<string>("warnings.burst_cd_data_header")}
        data={data}
        stat={(fa) => fa.burst_cd}
      />
    </DismissibleCallout>
  );
};

const SkillWarning = ({ data }: WarningProps) => {
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.skill_cd ?? false);

  return (
    <DismissibleCallout
        title={t<string>("warnings.skill_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.skill_body")}
      </p>
      <FailedActionDetails
          title={t<string>("warnings.skill_data_header")}
          data={data}
          stat={(fa) => fa.skill_cd} />
    </DismissibleCallout>
  );
};

const StaminaWarning = ({ data }: WarningProps) => {
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.insufficient_stamina ?? false);

  return (
    <DismissibleCallout
        title={t<string>("warnings.stamina_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.stamina_body")}
      </p>
      <FailedActionDetails
          title={t<string>("warnings.stamina_data_header")}
          data={data}
          stat={(fa) => fa.insufficient_stamina} />
    </DismissibleCallout>
  );
};

const SwapWarning = ({ data }: WarningProps) => {
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.swap_cd ?? false);

  return (
    <DismissibleCallout
        title={t<string>("warnings.swap_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.swap_body")}
      </p>
      <FailedActionDetails
          title={t<string>("warnings.swap_data_header")}
          data={data}
          stat={(fa) => fa.swap_cd} />
    </DismissibleCallout>
  );
};

const DashWarning = ({ data }: WarningProps) => {
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.statistics?.warnings?.dash_cd ?? false);

  return (
    <DismissibleCallout
        title={t<string>("warnings.dash_title")}
        intent={Intent.WARNING}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.dash_body")}
      </p>
      <FailedActionDetails
          title={t<string>("warnings.dash_data_header")}
          data={data}
          stat={(fa) => fa.dash_cd} />
    </DismissibleCallout>
  );
};

const IgnoreBurstEnergyMode = ({ data }: WarningProps) => {  
  const { t } = useTranslation();
  const [show, setShow] = useState(true);
  const visible = show && (data?.simulator_settings?.ignore_burst_energy ?? false);
  
  return (
    <DismissibleCallout
        title={t<string>("warnings.ignore_burst_energy_title")}
        intent={Intent.DANGER}
        show={visible}
        onDismiss={() => setShow(false)}>
      <p>
        {t<string>("warnings.ignore_burst_energy_body")}
      </p>
    </DismissibleCallout>
  );
};

type DetailsProps = {
  data: SimResults | null;
  title: string;
  stat: (x: FailedActions) => FloatStat | undefined;
}

const FailedActionDetails = ({ data, title, stat }: DetailsProps) => {
  const { i18n, t } = useTranslation();

  if (data?.character_details == null) {
    return null;
  }

  function fmt(val?: number) {
    return val?.toLocaleString(
        i18n.language, { maximumFractionDigits: 2 }) + t<string>("result.seconds_short");
  }

  const Item = ({ f, i }: { f: FloatStat | undefined, i: number }) => {
    if (f?.max == 0) {
      return null;
    }

    return (
      <>
        <div className="list-item">{t<string>("character_names." + data.character_details?.[i].name, { ns: "game" })}</div>
        <div>
          {
            `mean: ` + fmt(f?.mean)
            + ` | min: ` + fmt(f?.min)
            + ` | max: ` + fmt(f?.max)
            + ` | std: ` + fmt(f?.sd)
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
