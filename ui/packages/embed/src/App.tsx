import React from "react";
import { parsed } from ".";
import AvatarCard from "./AvatarCard";
import Details from "./Details";
import { ErrorBoundary } from "react-error-boundary";


function fallbackRender({ error  }) {
  return (
    <div id="card" role="alert">
      <p>Something went wrong:</p>
      <pre style={{ color: "red" }}>{error.message}</pre>
    </div>
  );
}

const App = ({ }) => {
  const [loaded, setLoaded] = React.useState(0);
  const [ready, setReady] = React.useState(false);

  if ("err" in parsed && parsed.err != "") {
    return <div id="card">{parsed.err}</div>;
  }

  const handleLoaded = () => {
    if (loaded + 1 == parsed.character_details?.length) {
      console.log("all loaded");
      setReady(true);
    }
    setLoaded(loaded + 1);
  };
  //draw some cards
  const cards = parsed.character_details?.map((c) => (
    <AvatarCard c={c} handleLoaded={handleLoaded} />
  ));

  return (
    <ErrorBoundary fallbackRender={fallbackRender}>
      <div className="bp4-dark flex flex-col align-middle justify-center h-full">
        <div
          id="card"
          className={ready ? "grid grid-cols-4 m-2" : "grid grid-cols-4 disabled"}
        >
          {cards}
        </div>
        <Details data={parsed} />
      </div>
    </ErrorBoundary>
  );
};

export default App;
