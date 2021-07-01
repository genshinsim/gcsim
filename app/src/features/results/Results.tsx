import { Card, Elevation, H4, H5 } from "@blueprintjs/core";
import { RootState } from "app/store";
import React from "react";
import { useSelector } from "react-redux";
import AvgModeResult from "./AvgModeResult";
import SingleModeResult from "./SingleModeResult";

function Results() {
  const { result, data } = useSelector((state: RootState) => {
    return {
      result: state.results.text,
      data: state.results.data,
    };
  });

  if (data === null) {
    return <div>No results. You need to run a simulation first.</div>;
  }

  return (
    <div>
      <div className="row">
        <div className="col-xs-10 col-xs-offset-1">
          <H4>Summary Result</H4>
          {"sim_duration" in data ? (
            <SingleModeResult data={data} />
          ) : (
            <AvgModeResult data={data} />
          )}

          <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
            <H5>Text Summary</H5>
            <pre>{result}</pre>
          </Card>
        </div>
      </div>
    </div>
  );
}

export default Results;
