import React from "react";
import { parsed } from ".";

const App = ({}) => {
  const [loaded, setLoaded] = React.useState(0);
  const [ready, setReady] = React.useState(false);

  const handleLoaded = () => {
    if (loaded + 1 == parsed.character_details?.length) {
      console.log("all loaded");
      setReady(true);
    }
    setLoaded(loaded + 1);
  };
  //draw some cards
  const cards = parsed.character_details?.map((c) => {
    return (
      <img
        key={c.name}
        src={`/api/assets/avatar/${c.name}.png`}
        onLoad={handleLoaded}
      />
    );
  });

  if (parsed.err) {
    return <div id="card">{parsed.err}</div>;
  }

  return (
    <div>
      <div
        id="card"
        className={ready ? "grid grid-cols-4" : "grid grid-cols-4 disabled"}
      >
        {cards}
      </div>
      <div>{`DPS: ${parsed.statistics.dps.mean}`}</div>
    </div>
  );
};

export default App;
