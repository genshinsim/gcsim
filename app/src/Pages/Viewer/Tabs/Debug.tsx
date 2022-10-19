import { Button, ButtonGroup, NonIdealState, Spinner, SpinnerSize } from "@blueprintjs/core";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { AdvancedPreset, AllDebugOptions, DebugPreset, DefaultDebugOptions, SimplePreset, VerbosePreset } from "~src/Components/Viewer/debugOptions";
import { Debugger } from "~src/Components/Viewer/DebugView";
import { Options } from "~src/Components/Viewer/Options";
import { DebugRow } from "~src/Components/Viewer/parse";
import { parseLogV2 } from "~src/Components/Viewer/parsev2";
import { SimResults } from "../SimResults";

const SAVED_DEBUG_KEY = "gcsim-debug-settings";

type Props = {
  data: SimResults | null;
  parsed: DebugRow[] | null;
  settingsState: [string[], (val: string[]) => void]; 
};

// TODO: The debugger should be refactored. This is a mess of passing around info
export default ({ data, parsed, settingsState }: Props) => {
  const [settings, setSettings] = settingsState;
  if (parsed == null || data?.character_details == null ) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  const names = data?.character_details?.map(c => c.name);
  return (
    <div className="flex flex-grow flex-col h-full gap-2 px-4">
      <Debugger data={parsed} team={names} searchable={{}} />
      <DebugOptions settings={settings} setSettings={setSettings} />
    </div>
  );
};

export function useDebugParser(data: SimResults | null, selected: string[]): DebugRow[] | null {
  return useMemo(() => {
    if (data?.initial_character == null || data.character_details == null || data?.debug == null) {
      return null;
    }

    return parseLogV2(
        data.initial_character,
        data?.character_details?.map(c => c.name),
        data.debug,
        selected);
  }, [data?.debug, data?.initial_character, data?.character_details, selected]);
}

export function useDebugSettings(): [string[], (val: string[]) => void] {
  const [selected, setSelected] = useState<string[]>(() => {
    const saved = localStorage.getItem(SAVED_DEBUG_KEY);
    if (saved) {
      const initialValue = JSON.parse(saved);
      return initialValue || DefaultDebugOptions;
    }
    return DefaultDebugOptions;
  });

  const setAndStore = (val: string[]) => {
    setSelected(val);
    localStorage.setItem(SAVED_DEBUG_KEY, JSON.stringify(val));
  };
  return [selected, setAndStore];
}

const DebugOptions = ({settings, setSettings}:
    {settings: string[], setSettings: (val: string[]) => void}) => {
  const { t } = useTranslation();
  const [isOpen, setOpen] = useState(false);

  const toggle = (t: string) => {
    const i = settings.indexOf(t);
      const next = [...settings];
      if (i === -1) {
        next.push(t);
      } else {
        next.splice(i, 1);
      }
      setSettings(next);
  };

  const presets = (opt: "simple" | "advanced" | "verbose" | "debug") => {
    switch (opt) {
      case "simple":
        setSettings(SimplePreset);
        return;
      case "advanced":
        setSettings(AdvancedPreset);
        return;
      case "verbose":
        setSettings(VerbosePreset);
        return;
      case "debug":
        setSettings(DebugPreset);
        return;
    }
  };

  return (
    <div className="w-full p-2">
      <ButtonGroup fill>
        <Button
            onClick={() => setOpen(true)}
            icon="cog"
            intent="primary"
            text={t("viewer.debug_settings")} />
      </ButtonGroup>
      <Options
          isOpen={isOpen}
          handleClose={() => setOpen(false)}
          handleClear={() => setSettings([])}
          handleResetDefault={() => setSettings(DefaultDebugOptions)}
          handleToggle={toggle}
          handleSetPresets={presets}
          selected={settings}
          options={AllDebugOptions} />
    </div>
  );
};