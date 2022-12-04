import { NonIdealState } from "@blueprintjs/core";
import qiqi from "./images/qiqi.png";
import kuki from "./images/kuki.png";
import ayaka from "./images/ayaka.png";
import xiao from "./images/xiao.png";

const images = [
  qiqi,
  kuki,
  xiao,
  ayaka,
];

// TODO: translation
export default ({}) => {
  return (
    <NonIdealState
        icon={<Icon />}
        title="Data not found"
        layout="horizontal" />
  );
};

let availableImages = [...images];

const Icon = ({}) => {
  const options = availableImages.length > 0 ? availableImages : [...images];
  const img = options.splice(Math.floor(Math.random() * options.length), 1);
  availableImages = options;
  return (
    <img src={img[0]} className="h-24" />
  );
};