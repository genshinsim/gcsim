import React from "react";
import { parsed } from ".";


const App = ({ }) => {
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
      <div className="card">
        <div className="char">
          <img
            key={c.name}
            src={`/api/assets/avatar/${c.name}.png`}
            onLoad={handleLoaded}
          />
        </div>
        <div className="equip">          
          <svg key={c.weapon.name} width="91" height="91">
            <filter id="outline">
              <feMorphology in="SourceAlpha" result="expanded" operator="dilate" radius="2.5"/>
              <feFlood flood-color="black"/>
              <feComposite in2="expanded" operator="in"/>
              <feGaussianBlur result="black" stdDeviation="1"/>
              <feMorphology in="SourceAlpha" result="expanded1" operator="dilate" radius="1"/>
              <feFlood flood-color="white"/>
              <feComposite in2="expanded1" operator="in" result="white"/>
              <feMerge>
                <feMergeNode in="black"/>
                <feMergeNode in="white"/>
                <feMergeNode in="SourceGraphic"/>
              </feMerge>
            </filter>
            <image filter="url(#outline)" href={`/api/assets/weapons/${c.weapon.name}.png`} height="85" width="85" x="3" y="3"/>
          </svg>
        </div>
      </div>

    );
  });

  if (parsed.err) {
    return <div id="card">{parsed.err}</div>;
  }
  const css = `
  .card {
    position: relative;
  }
  .equip {
    position: absolute;
    bottom: -10px;
    right: -3px;
    width: 91px;
  }
  `
  return (
    <div>
      <style>
        {css}
      </style>
      <div id="cards" className={ready ? "grid grid-cols-4" : "grid grid-cols-4 disabled"}>
        {cards}
      </div>
      <div>{`DPS: ${parsed.statistics.dps.mean}`}</div>
    </div>
    
  );  
};

export default App;
