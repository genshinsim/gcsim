import {
  Button,
  HTMLSelect,
  Intent,
  NonIdealState,
  OptionProps,
  Spinner,
  SpinnerSize,
} from "@blueprintjs/core";
import { Sample, SimResults } from "@gcsim/types";
import { useEffect, useMemo, useRef, useState } from "react";
import { DefaultSampleOptions, Sampler, SampleRow, parseLogV2 } from "../../Sample/Components";
import queryString from "query-string";
import { useTranslation } from "react-i18next";

const SAVED_SAMPLE_KEY = "gcsim-sample-settings";

type UseSampleData = {
  sample?: Sample;
  parsed: SampleRow[] | null;
  seed: string | null;
  searchable: { [key: number]: string[] };
  settings: string[];
  generating: boolean;
  setGenerating: (val: boolean) => void;
  setSample: (sample?: Sample) => void;
  setSettings: (val: string[]) => void;
  setSeed: (val: string | null) => void;
};

type Props = {
  sampler: (cfg: string, seed: string) => Promise<Sample>;
  data: SimResults | null;
  sample: UseSampleData;
  running: boolean;
};

// TODO: translation
// TODO: The sampler should be refactored. This is a mess of passing around info
export default ({ sampler, data, sample, running }: Props) => {
  const names = useMemo(() => {
    return data?.character_details?.map((c) => c.name);
  }, [data?.character_details]);

  if (names == null || data?.config_file == null || sample.generating) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  if (sample.sample == null || sample.parsed == null) {
    return (
      <NonIdealState
        icon="helper-management"
        action={<Generate sampler={sampler} data={data} sample={sample} running={running} />}
      />
    );
  }

  return (
    <div className="flex flex-grow flex-col gap-[15px] px-2">
      <Generate sampler={sampler} data={data} sample={sample} running={running} />
      <Sampler
          sample={sample.sample}
          data={sample.parsed}
          team={names}
          searchable={sample.searchable}
          settings={sample.settings}
          setSettings={sample.setSettings} />
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
  const { t } = useTranslation();
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
    { label: t<string>("viewer.seed_sample"), value: "sample" },
    // { label: "Random", value: "rand" },
    { label: t<string>("viewer.seed_min"), value: "min" },
    { label: t<string>("viewer.seed_max"), value: "max" },
    { label: t<string>("viewer.seed_p", { p: 25 }), value: "q1" },
    { label: t<string>("viewer.seed_p", { p: 50 }), value: "q2" },
    { label: t<string>("viewer.seed_p", { p: 75 }), value: "q3" },
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

    const parsed = queryString.parse(location.hash);
    parsed.sample = seed;
    location.hash = queryString.stringify(parsed);

    sample.setGenerating(true);
    sample.setSeed(seed);
    sampler(data.config_file ?? "", seed).then((out) => {
      sample.setSample(out);
      sample.setGenerating(false);
    });
  };

  return (
    <div className="flex flex-col gap-2 w-full mx-auto">
      <HTMLSelect
        options={options}
        value={value}
        onChange={(e) => setValue(e.currentTarget.value)}
        fill={true}
      />
      <Button
        large={true}
        text={t<string>("viewer.generate")}
        icon="refresh"
        intent={Intent.PRIMARY}
        disabled={disabled()}
        onClick={click}
        fill={true}
      />
    </div>
  );
};

export function useSample(
    running: boolean, data: SimResults | null, sampleOnLoad: boolean,
    sampler: (cfg: string, seed: string) => Promise<Sample>): UseSampleData {
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

  const [sample, SetSample] = useState<Sample | undefined>(undefined);
  const [generating, setGenerating] = useState(false);
  const [seed, setSeed] = useState<string | null>(null);
  
  // Special case where sim is rerunning. Want to reset any generated sample state
  useEffect(() => {
    if (running) {
      SetSample(undefined);
    }
  }, [running]);
  
  const initQuery = useRef(queryString.parse(location.hash));
  
  // if seed in url or sampleOnLoad is checked, load sample on viewer load
  useEffect(() => {
    const linkSeed = initQuery.current.sample as string;
    if ((sampleOnLoad || linkSeed) && sample == null && !generating && data?.config_file != null) {
      const seed = linkSeed ?? data.sample_seed;
      
      setGenerating(true);
      setSeed(seed);
      sampler(data.config_file ?? "", seed).then((out) => {
        SetSample(out);
        setGenerating(false);
      });
    }
  }, [data?.config_file, data?.sample_seed, generating, sample, sampleOnLoad, sampler]);

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
        sample.logs,
        selected);
  }, [sample, data?.initial_character, data?.character_details, selected]);

  const searchable = useMemo(() => {
    const out: { [key: number]: string[] } = {};
    if (parsed == null) {
      return out;
    }

    parsed.map((row, i) => {
      const results: string[] = [];
      row.slots.map((slot) => {
        slot.map((e) => {
          results.push(e.msg);
        });
      });
      out[i] = results;
    });
    return out;
  }, [parsed]);

  return useMemo(() => {
    return {
      sample: sample,
      parsed: parsed,
      seed: seed,
      searchable: searchable,
      settings: selected,
      generating: generating,
      setGenerating: setGenerating,
      setSample: SetSample,
      setSettings: setAndStore,
      setSeed: setSeed,
    };
  }, [generating, parsed, sample, searchable, seed, selected]);
}