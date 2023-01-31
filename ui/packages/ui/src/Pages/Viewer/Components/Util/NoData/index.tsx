import { NonIdealState } from "@blueprintjs/core";
import { memo, useRef } from "react";

import qiqi from "./images/qiqi.png";
import kuki from "./images/kuki.png";
import ayaka from "./images/ayaka.png";
import xiao from "./images/xiao.png";
import nahida from "./images/nahida.png";
import venti from "./images/venti.png";
import raiden from "./images/raiden.png";
import collei from "./images/collei.png";
import lisa from "./images/lisa.png";
import ayato from "./images/ayato.png";
import albedo from "./images/albedo.png";

const images = [
  qiqi,
  kuki,
  xiao,
  ayaka,
  nahida,
  venti,
  raiden,
  collei,
  lisa,
  ayato,
  albedo,
];

// TODO: translation
const NoData = ({}) => {
  return (
    <NonIdealState
        icon={<Icon />}
        title="Data not found"
        layout="horizontal" />
  );
};

let availableImages = [...images];

const Icon = ({}) => {
  const img = useRef<string | undefined>(image());
  return (
    <img src={img.current} className="h-24" />
  );
};

function image(): string {
  const options = availableImages.length > 0 ? availableImages : [...images];
  const img = options.splice(Math.floor(Math.random() * options.length), 1)[0];
  availableImages = options;
  return img;
}

export default memo(NoData);