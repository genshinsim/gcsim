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
    console.log(c.sets)
    var arr = Object.entries(c.sets)
    const sets = arr.filter(([key, value]) => typeof value === 'number' && value >=2) as [string, Number][]
    console.log(sets)
    var artifacts;
   if (sets.length >= 2){
      artifacts =  <g filter="url(#outlineb)">
        <image filter="url(#outlinew)" href={`/api/assets/artifacts/${sets[0][0]}_flower.png`} height="43" width="20.5" x="30" y="52" preserveAspectRatio="xMinYMid slice"></image>
        <image filter="url(#outlinew)" href={`/api/assets/artifacts/${sets[1][0]}_flower.png`} height="43" width="20.5" x="52.5" y="52" preserveAspectRatio="xMaxYMid slice"></image>
      </g>
    } else if (sets.length >= 1) {
      if (sets[0][1] >= 4) {
        artifacts = <image filter="url(#outlinew) url(#outlineb)" href={`/api/assets/artifacts/${sets[0][0]}_flower.png`} height="43" width="43" x="30" y="52"/>
      }else {
        artifacts = <image filter="url(#outlinew) url(#outlineb)" href={`/api/assets/artifacts/${sets[0][0]}_flower.png`} height="43" width="20.5" x="30" y="52" preserveAspectRatio="xMinYMid slice"></image>
      }
    } 
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
          <svg key={c.weapon.name} width="91" height="95">
            <filter id="outlinew">
              <feMorphology in="SourceAlpha" result="expanded" operator="dilate" radius="1"/>
              <feFlood flood-color="white"/>
              <feComposite in2="expanded" operator="in"/>
              <feComposite in="SourceGraphic"/>
            </filter>
            <filter id="outlineb">
              <feMorphology in="SourceAlpha" result="expanded" operator="dilate" radius="1.5"/>
              <feFlood flood-color="black"/>
              <feComposite in2="expanded" operator="in"/>
              <feGaussianBlur stdDeviation="1"/>
              <feComposite in="SourceGraphic"/>
            </filter>
            <image filter="url(#outlinew) url(#outlineb)" href={`/api/assets/weapons/${c.weapon.name}.png`} height="85" width="85" x="3" y="3"/>
            {artifacts}
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
    height: 140px;
  }
  .equip {
    position: absolute;
    bottom: -5px;
    right: -7px;
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
