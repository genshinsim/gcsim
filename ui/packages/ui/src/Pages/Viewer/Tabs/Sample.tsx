import {
  Button,
  ButtonGroup,
  HTMLSelect,
  Intent,
  NonIdealState,
  OptionProps,
  Spinner,
  SpinnerSize,
} from "@blueprintjs/core";
import { LogDetails, Sample, SimResults } from "@gcsim/types";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  AdvancedPreset,
  AllSampleOptions,
  DebugPreset,
  DefaultSampleOptions,
  SimplePreset,
  VerbosePreset,
  Sampler,
  Options,
  SampleRow,
  parseLogV2,
} from "../Components/Sample";

const SAVED_SAMPLE_KEY = "gcsim-sample-settings";

type UseSampleData = {
  logs?: LogDetails[];
  parsed: SampleRow[] | null;
  seed: string | null;
  settings: string[];
  generating: boolean;
  setGenerating: (val: boolean) => void;
  setLogs: (sample?: LogDetails[]) => void;
  setSettings: (val: string[]) => void;
  setSeed: (val: string | null) => void;
};

type Props = {
  sampler: (cfg: string, seed: string) => Promise<Sample>
  data: SimResults | null;
  sample: UseSampleData;
  running: boolean;
};

// TODO: translation
// TODO: The sampler should be refactored. This is a mess of passing around info
export default ({ sampler, data, sample, running }: Props) => {
  if (data?.character_details == null || data?.config_file == null || sample.generating) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  if (sample.parsed == null) {
    return (
      <NonIdealState
        icon="helper-management"
        action={<Generate sampler={sampler} data={data} sample={sample} running={running} />}
      />
    );
  }

  const msgs = useMemo(() => {
    const out: { [key: number]: string[] } = {};
    if (sample.parsed == null) {
      return out;
    }

    sample.parsed.map((row, i) => {
      const results: string[] = [];

      row.slots.map((slot) => {
        slot.map((e) => {
          results.push(e.msg);
        });
      });

      out[i] = results;
    });

    return out;
  }, [sample.parsed]);

  const names = data.character_details.map((c) => c.name);
  return (
    <div className="flex flex-grow flex-col h-full gap-2 px-4">
      <Generate sampler={sampler} data={data} sample={sample} running={running} />
      <Sampler data={sample.parsed} team={names} searchable={msgs} />
      <SampleOptions settings={sample.settings} setSettings={sample.setSettings} />
    </div>
  );
};

type GenerateProps = {
  sampler: (cfg: string, seed: string) => Promise<Sample>;
  data: SimResults;
  sample: UseSampleData;
  running: boolean;
}

const Generate = ({ sampler, data, sample, running }: GenerateProps) => {
  let startValue = "sample";
  switch (sample.seed) {
    case null:
      startValue = "sample";
      break;
    case data.sample_seed:
      startValue = "sample";
      break;
    case data.statistics?.min_seed:
      startValue = "min";
      break;
    case data.statistics?.max_seed:
      startValue = "max";
      break;
    case data.statistics?.p25_seed:
      startValue = "q1";
      break;
    case data.statistics?.p50_seed:
      startValue = "q2";
      break;
    case data.statistics?.p75_seed:
      startValue = "q3";
      break;
  }
  const [value, setValue] = useState(startValue);
  const options: OptionProps[] = [
    { label: "Sample Seed", value: "sample" },
    // { label: "Random", value: "rand" },
    { label: "Min Seed", value: "min" },
    { label: "Max Seed", value: "max" },
    { label: "P25 Seed", value: "q1" },
    { label: "P50 Seed", value: "q2" },
    { label: "P75 Seed", value: "q3" },
  ];

  const disabled = () => {
    return running && ["min", "max", "q1", "q2", "q3"].includes(value);
  };

  const click = () => {
    let seed = "0";
    switch (value) {
      case "sample":
        seed = data.sample_seed ?? seed;
        break;
      case "rand":
        seed = "" + Math.floor(Number.MAX_SAFE_INTEGER * Math.random());
        break;
      case "min":
        seed = data.statistics?.min_seed ?? seed;
        break;
      case "max":
        seed = data.statistics?.max_seed ?? seed;
        break;
      case "q1":
        seed = data.statistics?.p25_seed ?? seed;
        break;
      case "q2":
        seed = data.statistics?.p50_seed ?? seed;
        break;
      case "q3":
        seed = data.statistics?.p75_seed ?? seed;
        break;
    }

    sample.setGenerating(true);
    sample.setSeed(seed);
    sampler(data.config_file ?? "", seed).then((out) => {
      console.log(out);
      sample.setLogs(out.logs);
      sample.setGenerating(false);
    });
  };

  return (
    <>
      <HTMLSelect
        options={options}
        value={value}
        onChange={(e) => setValue(e.currentTarget.value)}
      />
      <Button
        large={true}
        text="Generate"
        icon="refresh"
        intent={Intent.PRIMARY}
        disabled={disabled()}
        onClick={click}
      />
    </>
  );
};

const SampleOptions = ({
  settings,
  setSettings,
}: {
  settings: string[];
  setSettings: (val: string[]) => void;
}) => {
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
    <div className="w-full p-2 pb-0">
      <ButtonGroup fill>
        <Button
          onClick={() => setOpen(true)}
          icon="cog"
          intent="primary"
          text={t<string>("viewer.sample_settings")}
        />
      </ButtonGroup>
      <Options
        isOpen={isOpen}
        handleClose={() => setOpen(false)}
        handleClear={() => setSettings([])}
        handleResetDefault={() => setSettings(DefaultSampleOptions)}
        handleToggle={toggle}
        handleSetPresets={presets}
        selected={settings}
        options={AllSampleOptions}
      />
    </div>
  );
};

export function useSample(running: boolean, data: SimResults | null): UseSampleData {
  const [selected, setSelected] = useState<string[]>(() => {
    const saved = localStorage.getItem(SAVED_SAMPLE_KEY);
    if (saved) {
      const initialValue = JSON.parse(saved);
      return initialValue || DefaultSampleOptions;
    }
    return DefaultSampleOptions;
  });

  const setAndStore = (val: string[]) => {
    setSelected(val);
    localStorage.setItem(SAVED_SAMPLE_KEY, JSON.stringify(val));
  };

  const [sample, SetSample] = useState<LogDetails[] | undefined>(undefined);
  const [generating, setGenerating] = useState(false);
  const [seed, setSeed] = useState<string | null>(null);

  // Special case where sim is rerunning. Want to reset any generated sample state
  useEffect(() => {
    if (running) {
      SetSample(undefined);
    }
  }, [running]);

  const parsed = useMemo(() => {
    if (data?.initial_character == null || data.character_details == null) {
      return null;
    }

    if (sample == null) {
      return null;
    }

    return parseLogV2(
        data.initial_character,
        data?.character_details?.map((c) => c.name),
        sample,
        selected);
  }, [sample, data?.initial_character, data?.character_details, selected]);

  return {
    logs: sample,
    parsed: parsed,
    seed: seed,
    settings: selected,
    generating: generating,
    setGenerating: setGenerating,
    setLogs: SetSample,
    setSettings: setAndStore,
    setSeed: setSeed,
  };
}