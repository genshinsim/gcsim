import React from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  TextArea,
  Intent,
  ButtonGroup,
  Button,
  Callout,
  Card,
  H4,
  FormGroup,
  Switch,
  NumericInput,
  Divider,
} from "@blueprintjs/core";
import { RootState } from "app/store";
import { runSim, setConfig } from "./simSlice";
import CharacterBuilder from "./CharacterBuilder";
import Import from "features/import/Import";
import ArtifactBuilder from "./ArtifactsBuilder";
import SampleConfig from "features/sample/SampleConfig";
import download from "downloadjs";
import dayjs from "dayjs";


function Dash() {
  const dispatch = useDispatch();
  const { config } = useSelector((state: RootState) => {
    return {
      config: state.sim.config,
    };
  });

  const [logDebug, setLogDebug] = React.useState<boolean>(true);
  const [duration, setDuration] = React.useState<number>(90);
  const [iter, setIter] = React.useState<number>(1000);
  const [workers, setWorkers] = React.useState<number>(24);

  const [openSample, setOpenSample] = React.useState<boolean>(false);
  const [openCharBuilder, setOpenCharBuilder] = React.useState<boolean>(false);
  const [openImport, setOpenImport] = React.useState<boolean>(false);
  const [openArtifactBuilder, setOpenArtifactBuilder] =
    React.useState<boolean>(false);

  React.useEffect(() => {
    // Update the document title using the browser API
    let config = localStorage.getItem("sim-config");
    if (config !== null) {
      dispatch(setConfig(config));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleRun = () => {
    dispatch(
      runSim({
        options: {
          log_details: true,
          duration: Math.round(duration),
          iter: Math.round(iter),
          debug: logDebug,
          workers: workers,
        },
        config: config,
      })
    );

  };

  const handleExport = () => {
    var now = dayjs();
    var filename =
      "gisim-config-" + now.format("YYYY-MM-DD-HH-mm-ssZ") + ".txt";
    download(config, filename, "text/plain");
  };

  const handleConfigChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    dispatch(setConfig(e.target.value));
  };

  return (
    <div>
      <div className="row">
        <div className="col-xs-9">
          <TextArea
            rows={30}
            fill
            large={true}
            intent={Intent.PRIMARY}
            onChange={handleConfigChange}
            value={config}
          ></TextArea>
          <br />
          <br />
          <ButtonGroup fill>
            <Button fill intent={Intent.PRIMARY} onClick={handleRun}>
              Run
            </Button>
          </ButtonGroup>

          <Callout
            intent="primary"
            style={{ marginBottom: "10px", marginTop: "10px" }}
            title={"Getting started"}
          >
            Get started by building a config file. Check out the{" "}
            <a
              href="https://github.com/srliao/gisim/wiki"
              target="_blank"
              rel="noreferrer"
            >
              {" "}
              wiki
            </a>{" "}
            for instructions. Or load one of the premade examples from the
            generator panel on the right.
            <br />
            <br />
            Note that your config below is saved automatically to your browser's
            local storage so that you will not lose it should you reload the
            page. However it is recommended that you export your config to a
            text file and save that somewhere safe in case you accidentally wipe
            your browser's local storage (or if you switch browsers)
          </Callout>
        </div>
        <div className="col-xs-3">
          <H4>Generator</H4>
          <Card>
            <Button
              fill
              style={{ marginTop: "10px", marginBottom: "10px" }}
              onClick={() => setOpenCharBuilder(true)}
            >
              Character
            </Button>
            <Button
              fill
              style={{ marginTop: "10px", marginBottom: "10px" }}
              onClick={() => setOpenArtifactBuilder(true)}
            >
              Artifacts
            </Button>
            <Button
              fill
              style={{ marginTop: "10px", marginBottom: "10px" }}
              onClick={() => setOpenImport(true)}
            >
              Import JSON
            </Button>

            <Button
              fill
              style={{ marginTop: "10px", marginBottom: "10px" }}
              onClick={() => setOpenSample(true)}
            >
              Load Premade Config
            </Button>

            <Button
              fill
              style={{ marginTop: "10px", marginBottom: "10px" }}
              onClick={handleExport}
            >
              Export Config
            </Button>
          </Card>
          <Divider style={{ marginTop: "10px", marginBottom: "10px" }} />

          <H4>Options</H4>

          <Card>
            <div>
              <FormGroup
                label="Iterations"
                helperText="number of iterations to run"
              >
                <NumericInput
                  value={iter}
                  onValueChange={(v) => setIter(v)}
                  min={1}
                />
              </FormGroup>
            </div>

            <FormGroup
              label="Duration"
              helperText="ignored if in damage mode as specified in the config"
            >
              <NumericInput
                value={duration}
                onValueChange={(v) => setDuration(v)}
                min={0}
              />
            </FormGroup>

            <FormGroup
              label="Workers"
              helperText="number of workers (threads) to use"
            >
              <NumericInput
                value={workers}
                onValueChange={(v) => setWorkers(v)}
                min={0}
              />
            </FormGroup>
            <FormGroup
              label="Debug"
              helperText="this option will capture the debug output of the last iteration"
            >
              <Switch
                checked={logDebug}
                onChange={(e) => setLogDebug(e.currentTarget.checked)}
              />
            </FormGroup>
          </Card>
        </div>
      </div>
      <br />
      <div className="row">
        <div className="col-xs-12"></div>
      </div>
      <SampleConfig
        isOpen={openSample}
        onClose={() => {
          setOpenSample(false);
        }}
      />
      <CharacterBuilder
        isOpen={openCharBuilder}
        onClose={() => {
          setOpenCharBuilder(false);
        }}
      />
      <ArtifactBuilder
        isOpen={openArtifactBuilder}
        onClose={() => {
          setOpenArtifactBuilder(false);
        }}
      />
      <Import
        isOpen={openImport}
        onClose={() => {
          setOpenImport(false);
        }}
      />
    </div>
  );
}

export default Dash;
