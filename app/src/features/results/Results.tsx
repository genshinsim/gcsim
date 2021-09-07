import { H4 } from "@blueprintjs/core";
import { RootState } from "app/store";
import React from "react";
import { useSelector } from "react-redux";
import AvgModeResult from "./AvgModeResult";

function Results() {
  const { data } = useSelector((state: RootState) => {
    return {
      // result: state.results.text,
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
          <AvgModeResult data={data} />
        </div>
      </div>
    </div>
  );
}

export default Results;
