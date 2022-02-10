import { Dialog } from "@blueprintjs/core";
import { Classes, Popover2, Tooltip2 } from "@blueprintjs/popover2";
import React from "react";
import { CharDetail } from "../DataType";
import weaponNames from "./weapons";

type weapMap = {
  [key: string]: string;
};

type Props = {
  char: CharDetail;
};

//stat types

const DEFP = 1;
const DEF = 2;
const HP = 3;
const HPP = 4;
const ATK = 5;
const ATKP = 6;
const ER = 7;
const EM = 8;
const CR = 9;
const CD = 10;
const Heal = 11;
const PyroP = 12;
const HydroP = 13;
const CryoP = 14;
const ElectroP = 15;
const AnemoP = 16;
const GeoP = 17;
const PhyP = 18;
const DendroP = 19;

function charBG(element: string) {
  switch (element) {
    case "cryo":
      return "bg-gradient-to-r from-gray-700 to-blue-300";
    case "hyro":
      return "bg-gradient-to-r from-gray-700 to-blue-500";
    case "pyro":
      return "bg-gradient-to-r from-gray-700 to-red-400";
    case "electro":
      return "bg-gradient-to-r from-gray-700 to-purple-300";
    case "anemo":
      return "bg-gradient-to-r from-gray-700 to-green-300";
    case "geo":
      return "bg-gradient-to-r from-gray-700 to-yellow-400";
  }
  return "bg-gray-700";
}

export default function Character({ char }: Props) {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);

  const arts: JSX.Element[] = [];

  for (const key in char.sets) {
    arts.push(
      <div className="w-8 flex flex-col rounded-md" key={key}>
        <Tooltip2
          content={<div className="bg-gray-800 m-2 p-2 rounded-md">{key}</div>}
        >
          <img
            key="key"
            src={`/images/artifacts/${key}_flower.png`}
            alt={key}
            className="w-full h-auto "
          />
        </Tooltip2>

        <span className="text-center text-xs">{char.sets[key]}</span>
      </div>
    );
  }

  return (
    <div className="min-h-48 bg-gray-600 shadow rounded-md text-sm m-2 flex flex-col justify-center gap-2">
      <div
        className={
          "character-parent flex flex-row pt-4 pl-4 pr-2 -mt-2 rounded-t-md " +
          charBG(char.element)
        }
      >
        <div className="character-header rounded-t-md" />
        <div className="character-name font-medium m-4 capitalize">
          {char.name}
        </div>
        <div className="w-1/2">
          <div className="rounded-md pl-1 pr-1 mt-6">
            <div>
              Level {char.level}/{char.max_level}
            </div>
            <div>Cons {char.cons}</div>
            <div>
              Talent: {char.talents.attack}/{char.talents.skill}/
              {char.talents.burst}
            </div>
            <div className="mt-2 mr-2 grid grid-cols-5">{arts}</div>
          </div>
        </div>
        <div className="w-1/2">
          <img
            src={"/images/avatar/" + char.name + ".png"}
            alt={char.name}
            className="w-full h-auto "
          />
        </div>
      </div>

      <div className="weapon-parent ml-2 mr-2 mb-2 p-2 bg-gray-800 rounded-md">
        <div className="flex flex-row">
          <div className="rounded-md">
            <img
              src={`/images/weapons/${char.weapon.name}.png`}
              alt={char.weapon.name}
              className="object-contain h-16"
            />
          </div>
          <div className="flex-grow">
            <div className="font-medium ml-2 text-left">
              {weaponNames[char.weapon.name]}
            </div>

            <div className="ml-2 justify-center items-center rounded-md">
              <div>Level {char.weapon.level}/90</div>
              <div>Refinement {char.weapon.refine}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function SubstatView({ char }: Props) {
  return (
    <div className="flex flex-col m-2 p-2">
      <div className="flex flex-row">
        <div className="font-bold">Substats</div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="hp_primary"
              className="svg-inline--fa fa-hp_primary fa-w-13 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 11.104 14"
            >
              <path
                fill="currentColor"
                d="M2.052,7.654A.978.978,0,0,1,2.5,7.083c1.511-.85,3.586,2.117,6.544.548C10.971,13.714.151,13.878,2.052,7.654ZM5.552,14C2.177,14-1.2,11.579.416,6.74A18.543,18.543,0,0,1,5.121.213.748.748,0,0,1,5.552,0a.751.751,0,0,1,.431.212A18.543,18.543,0,0,1,10.688,6.74C12.3,11.579,8.927,14,5.552,14Zm.22-12.19a.639.639,0,0,0-.22-.075h0a.649.649,0,0,0-.221.075c-1.71,1.324-4.08,5.282-3.941,7.4a4.019,4.019,0,0,0,4.162,3.543h0A4.019,4.019,0,0,0,9.714,9.215C9.853,7.092,7.483,3.134,5.772,1.81Z"
              />
            </svg>
          </span>
          <span>HP</span>
        </div>
        <div className="flex-grow text-right">
          {(char.stats[HPP] * 100).toFixed(2) + "%"} | {char.stats[HP]}
        </div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="Atk"
              className="svg-inline--fa fa-Atk fa-w-17 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 14 13.977"
            >
              <path
                fill="currentColor"
                d="M7.78,1.621,9.423,3.264l1.928-1.927L11.326.219A.228.228,0,0,1,11.554,0h2.2a.228.228,0,0,1,.228.228c-.121,2.661.556,2.457-1.337,2.4L10.712,4.553,12.356,6.2a.228.228,0,0,1,0,.322c-1.167,1.208-.775.907-1.892-.106L3.313,13.564A.457.457,0,0,1,3,13.7a21.32,21.32,0,0,1-2.954.239,21.172,21.172,0,0,1,.238-2.954.451.451,0,0,1,.134-.318L7.564,3.513l-.838-.838a.229.229,0,0,1,0-.323l.732-.731A.228.228,0,0,1,7.78,1.621Z"
              />
            </svg>
          </span>
          <span>Attack</span>
        </div>
        <div className="flex-grow text-right">
          {(char.stats[ATKP] * 100).toFixed(2) + "%"} | {char.stats[ATK]}
        </div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="Def"
              className="svg-inline--fa fa-Def fa-w-15 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 12.883 14"
            >
              <path
                fill="currentColor"
                d="M12.883.726a.291.291,0,0,0-.175-.268C12.3.286,10.944,0,6.442,0S.584.286.176.458A.291.291,0,0,0,0,.727v7.44a.868.868,0,0,0,.125.453c1.579,2.6,5.347,4.855,6.161,5.339a.292.292,0,0,0,.3,0c.789-.482,4.559-2.688,6.168-5.335a.868.868,0,0,0,.127-.455ZM6.441,11.968C6.5,11.981,2.882,9.951,1.617,7.6V1.565s0-.452,4.824-.452Z"
              />
            </svg>
          </span>
          <span>Defense</span>
        </div>
        <div className="flex-grow text-right">
          {(char.stats[DEFP] * 100).toFixed(2) + "%"} | {char.stats[DEF]}
        </div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="ele_mas"
              className="svg-inline--fa fa-ele_mas fa-w-18 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 14 13.113"
            >
              <path
                fill="currentColor"
                d="M8.076,7.71l-.017-.049A4.335,4.335,0,0,0,7.3,6.353a4.431,4.431,0,0,0-.325-.346A2.113,2.113,0,1,0,7,1.781a2.144,2.144,0,0,0-1.838,3.18,4.374,4.374,0,0,0-1.2-.168,4.42,4.42,0,0,0-.755.066l-.038.007C1.836-.682,10.7-2.114,10.962,3.9A3.985,3.985,0,0,1,8.076,7.71Zm3.662-2.137a3.949,3.949,0,0,0-.626-.235,4.473,4.473,0,0,1-1.105,1.7h.031A2.113,2.113,0,1,1,7.925,9.151,4.09,4.09,0,0,0,7.9,8.706,3.968,3.968,0,0,0,6.037,5.775l-.19-.11A3.963,3.963,0,1,0,6.492,12.2c.082-.068.16-.14.236-.214L6.7,11.949a4.367,4.367,0,0,1-.891-1.765A2.112,2.112,0,1,1,4.926,7.27q.1.051.189.111a2.111,2.111,0,0,1,.942,1.49,2.159,2.159,0,0,1,.018.28,3.963,3.963,0,1,0,5.663-3.578Z"
              />
            </svg>
          </span>
          <span>EM</span>
        </div>
        <div className="flex-grow text-right">{char.stats[EM]}</div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="EnergyRecharge_primary"
              className="svg-inline--fa fa-EnergyRecharge_primary fa-w-17 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 14 13.282"
            >
              <path
                fill="currentColor"
                d="M14,6.641A6.642,6.642,0,0,1,.928,8.3h0L0,8.737.961,4.8l.012.012L2.392,6.3l1.37,1.433-1.23.143A4.981,4.981,0,1,0,7.359,1.66V0A6.641,6.641,0,0,1,14,6.641Z"
              />
            </svg>
          </span>
          <span>ER</span>
        </div>
        <div className="flex-grow text-right">
          {char.stats[ER].toFixed(2) + "%"}
        </div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="faga"
              data-icon="CritRate"
              className="svg-inline--fa fa-CritRate fa-w-16 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 14 14"
            >
              <path
                fill="currentColor"
                d="M0,14,3.5,6.236,1.973,1.464,7.255,3.5,14,0,10.5,7.255l1.527,4.772L7.255,10.5,0,14"
              />
            </svg>
          </span>
          <span>CR</span>
        </div>
        <div className="flex-grow text-right">
          {(char.stats[CR] * 100).toFixed(2) + "%"}
        </div>
      </div>
      <div className="flex flex-row ml-2">
        <div>
          <span className="substat-header">
            <svg
              aria-hidden="true"
              focusable="false"
              data-prefix="fas"
              data-icon="dice-d20"
              className="svg-inline--fa fa-dice-d20 fa-w-15 "
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 480 512"
            >
              <path
                fill="currentColor"
                d="M106.75 215.06L1.2 370.95c-3.08 5 .1 11.5 5.93 12.14l208.26 22.07-108.64-190.1zM7.41 315.43L82.7 193.08 6.06 147.1c-2.67-1.6-6.06.32-6.06 3.43v162.81c0 4.03 5.29 5.53 7.41 2.09zM18.25 423.6l194.4 87.66c5.3 2.45 11.35-1.43 11.35-7.26v-65.67l-203.55-22.3c-4.45-.5-6.23 5.59-2.2 7.57zm81.22-257.78L179.4 22.88c4.34-7.06-3.59-15.25-10.78-11.14L17.81 110.35c-2.47 1.62-2.39 5.26.13 6.78l81.53 48.69zM240 176h109.21L253.63 7.62C250.5 2.54 245.25 0 240 0s-10.5 2.54-13.63 7.62L130.79 176H240zm233.94-28.9l-76.64 45.99 75.29 122.35c2.11 3.44 7.41 1.94 7.41-2.1V150.53c0-3.11-3.39-5.03-6.06-3.43zm-93.41 18.72l81.53-48.7c2.53-1.52 2.6-5.16.13-6.78l-150.81-98.6c-7.19-4.11-15.12 4.08-10.78 11.14l79.93 142.94zm79.02 250.21L256 438.32v65.67c0 5.84 6.05 9.71 11.35 7.26l194.4-87.66c4.03-1.97 2.25-8.06-2.2-7.56zm-86.3-200.97l-108.63 190.1 208.26-22.07c5.83-.65 9.01-7.14 5.93-12.14L373.25 215.06zM240 208H139.57L240 383.75 340.43 208H240z"
              />
            </svg>
          </span>
          <span>CD</span>
        </div>
        <div className="flex-grow text-right">
          {(char.stats[CD] * 100).toFixed(2) + "%"}
        </div>
      </div>
    </div>
  );
}
