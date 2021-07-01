import React from "react";

import Dash from "features/sim/Dash";

function Sim() {
  return (
    <div style={{ marginLeft: "20px", marginRight: "20px" }}>
      <div className="row">
        <div className="col-xs-10 col-xs-offset-1">
          <Dash />
        </div>
      </div>
    </div>
  );
}

export default Sim;
