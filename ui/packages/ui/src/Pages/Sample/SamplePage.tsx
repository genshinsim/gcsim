import { Alert, ButtonGroup, Intent, NonIdealState, Position, Spinner, SpinnerSize, Toaster } from "@blueprintjs/core";
import { Sample } from "@gcsim/types";
import { useMemo, useRef, useState } from "react";
import { useHistory } from "react-router";
import { CopyToClipboard, SendToSimulator } from "../../Components/Buttons";
import { CharacterCard } from "../../Components/Cards";
import { characterCardsClassNames } from "../Viewer/Components/Overview/TeamHeader";
import { DefaultSampleOptions, parseLogV2, Sampler, SampleRow } from "./Components";

const SAVED_SAMPLE_KEY = "gcsim-sample-settings";

type UseSampleData = {
  parsed: SampleRow[] | null;
  team?: string[];
  searchable: { [key: number]: string[] };
  settings: string[];
  setSettings: (val: string[]) => void;
};

type Props = {
  sample: Sample | null;
  error: string | null;
  retry?: () => void;
}

export default ({ sample, error, retry }: Props) => {
  const data = useSample(sample);
  const copyToast = useRef<Toaster>(null);

  if (sample == null || data.team == null || data.parsed == null) {
    return (
      <>
        <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />
        <ErrorAlert msg={error} retry={retry} />
      </>
    );
  }

  const cardClass = characterCardsClassNames(sample.character_details?.length ?? 4);
  return (
    <div className="flex flex-col gap-2 w-full 2xl:mx-auto 2xl:container py-6">
      <div className="flex flex-row justify-between px-6 pb-2">
        <span className="text-lg font-bold font-mono">
          {"Targets: " + sample.target_details?.length}
        </span>
        <ButtonGroup>
          <CopyToClipboard
              copyToast={copyToast}
              config={sample.config}
              className="hidden ml-[7px] sm:flex" />
          <SendToSimulator config={sample.config} />
        </ButtonGroup>
      </div>
      <div className="flex flex-row gap-2 justify-center flex-wrap">
        {sample.character_details?.map((c) => (
          <CharacterCard
              key={c.name}
              char={c}
              showDetails={false}
              stats={[]}
              statsRows={0}
              className={cardClass} />
        ))}
      </div>
      <div className="flex flex-grow flex-col gap-[15px] px-4">
        <Sampler
            sample={sample}
            data={data.parsed}
            team={data.team}
            searchable={data.searchable}
            settings={data.settings}
            setSettings={data.setSettings} />
        <ErrorAlert msg={error} retry={retry} />
      </div>
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
    </div>
  );
};

type ErrorProps = {
  msg: string | null;
  retry?: () => void;
}

const ErrorAlert = ({ msg, retry }: ErrorProps) => {
  const history = useHistory();

  let cancelButtonText: string | undefined;
  let onCancel: (() => void) | undefined;
  if (retry != null) {
    cancelButtonText = "Retry";
    onCancel = () => retry();
  }

  return (
    <Alert
        isOpen={msg != null}
        onConfirm={() => history.push("/")}
        onCancel={onCancel}
        canEscapeKeyCancel={false}
        canOutsideClickCancel={false}
        confirmButtonText="Close"
        cancelButtonText={cancelButtonText}
        intent={Intent.DANGER}>
      <p>{msg}</p>
    </Alert>
  );
};


function useSample(sample: Sample | null): UseSampleData {
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
  
  const parsed = useMemo(() => {
    if (sample?.initial_character == null || sample?.character_details == null) {
      return null;
    }

    return parseLogV2(
        sample.initial_character,
        sample.character_details.map((c) => c.name),
        sample.logs,
        selected);
  }, [sample?.character_details, sample?.initial_character, sample?.logs, selected]);

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

  return {
    parsed: parsed,
    team: sample?.character_details?.map((c) => c.name),
    searchable: searchable,
    settings: selected,
    setSettings: setAndStore
  };
}