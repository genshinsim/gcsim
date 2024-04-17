import { Button } from "@blueprintjs/core";
import tanuki from "images/tanuki.png";
import React from "react";
import { Trans, useTranslation } from "react-i18next";

export function Warning() {
  const { t } = useTranslation();
  const [hide, setHide] = React.useState<boolean>((): boolean => {
    return localStorage.getItem("hide-warning") === "true";
  });
  React.useEffect(() => {
    localStorage.setItem("hide-warning", hide.toString());
  }, [hide]);

  if (hide) {
    return (
      <div className="flex flex-col py-0 min-[1300px]:w-[970px]">
        <div className="ml-auto">
          <Button
            small
            intent="success"
            onClick={() => {
              setHide(false);
            }}
          >
            {t<string>("db.readme_show")}
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="relative flex flex-col gap-2 items-center bg-slate-900 px-5 py-0 border border-blue-800 min-[1300px]:w-[970px] max-w-[970px]">
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
          {t<string>("db.readme_header")}
        </div>
        <img src={tanuki} className="w-15 h-10 mx-0" />
      </div>
      <div className="space-y-3 pb-3 text-s leading-5 text-gray-400">
        <Trans i18nKey="db.readme_body">
          <p />
          <p>{{ rerun: t<string>("viewer.rerun") }}</p>
          <p className="font-semibold leading-6 text-gray-200" />
        </Trans>
      </div>
    </div>
  );
}
