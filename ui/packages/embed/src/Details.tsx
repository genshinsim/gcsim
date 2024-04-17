//TODO(kyle): data should be typed model.SimulationResult
const Details = ({ data }: { data: any }) => {
  let valid = "key_type" in data && data.key_type !== "";
  if ("modified" in data) {
    valid = valid && data.modified === false;
  }

  return (
    <div className="flex flex-col  grow">
      <div className="flex flex-row ml-2 mr-2 ">
        <div>Created on {data.build_date}</div>
        {valid ? (
          <div className="ml-auto  text-green-700 font-semibold">
            Validated with {data.key_type} key
          </div>
        ) : (
          <div className="ml-auto  text-red-700 font-semibold">
            Warning: this sim may be modified
          </div>
        )}
      </div>
      {!("schema_version" in data) ? (
        <div className="flex flex-col ml-2 mr-2 mb-2 mt-2 p-2 rounded bg-[#252A31] ">
          Simulation out dated; Please rerun
        </div>
      ) : (
        <div className="flex flex-col ml-2 mr-2 mb-2 mt-2 p-2 rounded bg-[#252A31] ">
          <div className="grid grid-cols-3">
            <div>Mode: {data.mode == 1 ? "Dur" : "TTK"} </div>
            <div>{`Sim duration: ${Math.round(
              data.statistics.duration.mean
            )}s`}</div>
            <div>Iters: {data.statistics.iterations} </div>
          </div>
          <div className="grid grid-cols-3">
            <div>{`Target count: ${data.target_details.length}`}</div>
            <div>{`DPS: ${Math.round(
              data.statistics.dps.mean
            ).toLocaleString()}`}</div>
            <div>{`DPS/Target: ${Math.round(
              data.statistics.dps.mean / data.target_details.length
            ).toLocaleString()}`}</div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Details;
