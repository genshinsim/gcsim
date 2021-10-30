import React from "react";
import MainImage from "../images/main.png";
import DebuggerImage from "../images/debugger.png";
import ImporterImage from "../images/importer.png";
import ActionImage from "../images/actions.png";
import ResultsImage from "../images/results.png";
import LeftArrow from "../arrow_back_ios_white_48dp.svg";
import RightArrow from "../arrow_forward_ios_white_48dp.svg";
import { useLocation } from "wouter";

export default function Home() {
  const [pos, setPos] = React.useState<number>(0);
  const [_, setLocation] = useLocation();

  const images = [
    MainImage,
    ActionImage,
    ImporterImage,
    ResultsImage,
    DebuggerImage,
  ];

  const handleLeft = () => {
    let next = pos - 1;
    if (next < 0) {
      next = images.length - 1;
    }
    setPos(next);
  };
  const handleRight = () => {
    let next = pos + 1;
    if (next == images.length) {
      next = 0;
    }
    setPos(next);
  };
  const titleString = [
    "gcsim is a Monte Carlo combat simulator for Genshin Impact.",
    "It runs on a fully configurable action list",
    "It can import from Genshin Optimizer to save you time",
    "Results are shown with detailed breakdowns",
    "The debugger can show you exactly what the sim is doing",
  ];

  return (
    <div className="flex-grow flex flex-col p-10 text-center items-center">
      <p className="text-2xl md:text-xl font-medium p-4 mb-2">
        {titleString[pos]}
      </p>
      <div className="flex flex-row items-center">
        <div className="ml-auto p-2">
          <img
            src={LeftArrow}
            alt="prev"
            className="p-2 rounded-md hover:bg-gray-600"
            onClick={handleLeft}
          />
        </div>
        <div>
          <img
            src={images[pos]}
            alt="image"
            style={{
              width: "50vw",
            }}
          />
        </div>
        <div className="mr-auto p-2">
          <img
            src={RightArrow}
            alt="prev"
            className="p-2 rounded-md hover:bg-gray-600"
            onClick={handleRight}
          />
        </div>
      </div>
      <div className="grid md:grid-cols-2 gap-2 md:w-3/4 sm:grid-cols-1 sm:w-full mt-4">
        <div
          className="p-4 rounded-md bg-blue-800 hover:bg-blue-700 cursor-pointer"
          onClick={() => {
            setLocation("/getting-started");
          }}
        >
          Get Started
        </div>
        <div
          className="p-4 rounded-md bg-green-800 hover:bg-green-700 cursor-pointer"
          onClick={() => {
            setLocation("/db");
          }}
        >
          Find an Action List
        </div>
      </div>
    </div>
  );
}
