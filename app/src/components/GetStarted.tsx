import React from "react";
import MainImage from "../images/main.png";
import DebuggerImage from "../images/debugger.png";
import ImporterImage from "../images/importer.png";
import ActionImage from "../images/actions.png";
import ResultsImage from "../images/results.png";
import TutorialImage from "../images/tutorial.gif";
import { Link } from "wouter";

export default function GetStarted() {
  const [pos, setPos] = React.useState<number>(0);

  const images = [
    ActionImage,
    MainImage,
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

  return (
    <div className="flex-grow flex flex-col p-10 w-2/3 mx-auto">
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2">
        Introduction
      </p>
      <p className="p-2">
        gcsim is a Monte Carlo combat simulator for Genshin Impact. It attempts
        to simulate (to the best of our knowledge and ability) combat against a
        dummy target over a user specified time frame. If you are a World of
        Warcraft player and you have used SimulationCraft then this app should
        feel right at home.
      </p>
      <p className="p-2">
        Get started by downloading the desktop application on our GitHub release
        page:
      </p>
      <a
        href="https://github.com/genshinsim/gsimui/releases"
        target="_blank"
        className="p-4 rounded-md bg-blue-800 hover:bg-blue-700 mx-auto mt-3 "
      >
        Download Here
      </a>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2 mt-4">
        Just show me how to use this thing!!!
      </p>
      <div className="p-6">
        <img
          src={TutorialImage}
          alt="tutorial"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2  mt-4">
        How does it work
      </p>
      <p className="p-2">
        The simulator runs on a text based config file (very similar to WoW
        SimC). The config file can be broken down into two parts: the team
        configuration and the actions configuration.
      </p>
      <div className="p-6">
        <img
          src={ActionImage}
          alt="image"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2 mt-4">
        Team configuration
      </p>
      <p className="p-2">
        To make life easier, the app has a built in team configuration that let
        you pick which character you wish to use as well as their weapons,
        artifacts etc... This is combined with the action list to form the
        config file used by the simulator's core engine
      </p>
      <div className="p-6">
        <img
          src={MainImage}
          alt="image"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="p-2">
        Alternatively, you can also import your character data from{" "}
        <a
          href="https://frzyc.github.io/genshin-optimizer/#/"
          target="_blank"
          className="text-blue-600 hover:underline"
        >
          Genshin Optimizer
        </a>{" "}
        as well as any inventory scanner that supports the{" "}
        <a
          href="https://frzyc.github.io/genshin-optimizer/#/doc"
          className="text-blue-600 hover:underline"
          target="_blank"
        >
          GOOD
        </a>{" "}
        format (Currently to the best of our knowledge, just{" "}
        <a
          href="https://github.com/Andrewthe13th/Genshin_Scanner/releases"
          className="text-blue-600 hover:underline"
          target="_blank"
        >
          this one
        </a>
        ).
      </p>
      <div className="p-6">
        <img
          src={ImporterImage}
          alt="image"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2 mt-4">
        Actions configuration
      </p>
      <p className="p-2">
        The actions configuration (or action lists) are basically instructions
        that tells the simulator how to play the team you have specified. This
        is the bread and butter of the simulator.
      </p>
      <p className="p-2">
        The action list comprises of a list of actions in priority sequence that
        the simulator will execute. This is important as the actions are not
        listed in <strong>sequential</strong> order, but rather in{" "}
        <strong>priority</strong>. To put it simply, the simulator decides which
        action to execute based on the first action on the list that has is
        ready and has it's conditions fulfiled every time it needs something to
        execute. This allows for much more flexibility and let's the simulator
        decides conditionally on what to use next, similar to how an actual
        player may play. For anyone coming from WoW SimC, this should be no
        stranger to you.
      </p>
      <p className="p-2">
        For details on how to <strong>write</strong> action lists, head over to
        the wiki on our GitHub repo{" "}
        <a
          href="https://github.com/genshinsim/gsim/wiki"
          className="text-blue-600 hover:underline"
          target="_blank"
        >
          here
        </a>
        . However, we recognize that not everyone wants to learn/write action
        list. Some of us just wants to plug in our characters and see big
        numbers! To help with that, we have on this website a collection of
        action lists created by others
      </p>
      <Link href="/db">
        <a
          href="/db"
          className="p-4 rounded-md bg-green-800 hover:bg-green-700 mx-auto mt-3 "
        >
          Check out our action list DB
        </a>
      </Link>
      <p className="p-2">
        Simply copy and paste your favorite action list directly into the{" "}
        <strong>Actions</strong> box on the simulator's main screen! But make
        sure that your team of characters matches first!
      </p>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2 mt-4">
        Results and Debugging
      </p>
      <p className="p-2">
        Simulation results are displayed with detailed break down of the team's
        damage composition, along with other statistic to help you analyze the
        team's performance
      </p>
      <div className="p-6">
        <img
          src={ResultsImage}
          alt="image"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="p-2">
        The simulator also has built in debugging tool to help you see the
        actions executed by the sim and all related calculations. Make sure the{" "}
        <strong>debug</strong> option is selected in the sim options and
        checkout the Debug tab after you have ran a simulation.
      </p>
      <div className="p-6">
        <img
          src={DebuggerImage}
          alt="image"
          style={{
            width: "50vw",
          }}
        />
      </div>
      <p className="p-2">
        The debug view shows you chronogically (laid out vertically) what is
        happening frame by frame. The left most column shows you the current
        frame number. The other 5 columns shows you the actions taken by the
        simulator (for non character specific events) and by each character.
      </p>
      <p className="p-2">
        Character that are currently active for that frame is highlighted in
        light gray.
      </p>
      <p className="p-2">
        To see details for each action, simply click on any of the tags.
        Clicking on options will allow you to show or hide additional tags. By
        default, not all the tags are shown (or else there would be too much
        information!). So be sure to check out what other information there are.
      </p>
      <p className="lg:text-4xl md:text-3xl sm:text-xl font-bold mb-2 mt-4">
        Developers
      </p>
      <p className="p-2">
        The simulator is still under heavy development. Not all of the
        characters have been implemented (working very hard here). If you wish
        to help out pop in our discord and shoot us a message!
      </p>
      <p className="p-2">
        The app is made up of two primary packages: the{" "}
        <a
          href="https://github.com/genshinsim/gsim"
          className="text-blue-600 hover:underline"
          target="_blank"
        >
          core
        </a>{" "}
        and the{" "}
        <a
          href="https://github.com/genshinsim/gsimui"
          className="text-blue-600 hover:underline"
          target="_blank"
        >
          UI
        </a>
        . The core is the engine of the simulator, containing all the code to
        perform the actual simulation. It includes various command line tools
        for actually running the simulator. The UI is actually just a wrapper
        around the command line tool.
      </p>
    </div>
  );
}
