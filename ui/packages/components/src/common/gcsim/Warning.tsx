import {Button} from '@blueprintjs/core';
import tanuki from '../../images/tanuki.png';
import React from 'react';
import {Trans, useTranslation} from 'react-i18next';

interface WarningProps {
  hideKey: string;
  headerKey: string;
  bodyKey: string;
  bodyComponents?: Record<string, React.ReactElement>;
  showButton?: boolean;
  className?: string;
}

export function Warning({
  hideKey,
  headerKey,
  bodyKey,
  bodyComponents,
  showButton = true,
  className = "bg-slate-900 border-blue-800",
}: WarningProps) {
  const { t } = useTranslation();
  const [hide, setHide] = React.useState<boolean>(() => {
    return localStorage.getItem(hideKey) === 'true';
  });
  React.useEffect(() => {
    localStorage.setItem(hideKey, hide.toString());
  }, [hide, hideKey]);

  if (hide) {
    if (!showButton) {
      return <></>;
    }
    return (
      <div className="flex flex-col py-0 min-[1300px]:w-[1100px]">
        <div className="ml-auto">
          <Button
            small
            intent="success"
            onClick={() => {
              setHide(false);
            }}>
            {t('db.readme_show')}
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative flex flex-col gap-2 items-center px-5 py-0 border min-[1300px]:w-[1100px] ${className}`}>
      <div className="absolute top-1 right-1">
        <Button
          icon="cross"
          small
          intent="danger"
          onClick={() => {
            setHide(true);
          }}
        />
      </div>
      <div className="inline-flex pt-4">
        <img src={tanuki} className="w-15 h-10 mx-0" />
        <div className="font-semibold px-3 pt-2 text-xl w-50 text-gray-200">
          {t(headerKey)}
        </div>
        <img src={tanuki} className="w-15 h-10 mx-0" />
      </div>
      <div className="space-y-3 pb-3 text-s leading-5 text-gray-400">
        <Trans i18nKey={bodyKey} components={bodyComponents}>
          <p />
          <p>{{ rerun: t('viewer.rerun') }}</p>
          <p className="font-semibold leading-6 text-gray-200" />
        </Trans>
      </div>
    </div>
  );
}

export const RiskWarning = () => (
  <Warning
    hideKey="hide-warning-risk"
    headerKey="warnings.gcsim_risk_title"
    bodyKey="warnings.gcsim_risk_body"
    bodyComponents={{
      b: <b />,
      p: <p className="text-gray-200" />,
      discordlink: <a href="https://discord.com/channels/845087716541595668/983391844631212112/1413229755687505951" target="_blank" rel="noreferrer" />,
      githublink: <a href="https://github.com/genshinsim/gcsim" target="_blank" rel="noreferrer" />,
      issueslink: <a href="https://github.com/genshinsim/gcsim/issues" target="_blank" rel="noreferrer" />,
    }}
    className="bg-red-950 border-red-800"
    showButton={false}
  />
);