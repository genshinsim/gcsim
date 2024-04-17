import { memo, useRef } from "react";
import { useTranslation } from "react-i18next";
import { cn } from "../../lib/utils";

const images = [
  "/api/assets/misc/qiqi.png",
  "/api/assets/misc/kuki.png",
  "/api/assets/misc/ayaka.png",
  "/api/assets/misc/xiao.png",
  "/api/assets/misc/nahida.png",
  "/api/assets/misc/venti.png",
  "/api/assets/misc/raiden.png",
  "/api/assets/misc/collei.png",
  "/api/assets/misc/lisa.png",
  "/api/assets/misc/ayato.png",
  "/api/assets/misc/albedo.png",
];

type Props = {
  className?: string;
};

// TODO: translation
export const NoData = memo(({ className = "h-24" }: Props): JSX.Element => {
  const { t } = useTranslation();
  const cc = cn(
    "flex flex-row items-center font-bold text-gray-400 gap-5 text-lg",
    className
  );
  return (
    <div className={cc}>
      <NoDataIcon className={className} />
      <div>{t("common.data_not_found")}</div>
    </div>
  );
});

let availableImages = [...images];

export const NoDataIcon = ({ className }: Props) => {
  const img = useRef<string | undefined>(image());
  return <img src={img.current} className={className} />;
};

function image(): string {
  const options = availableImages.length > 0 ? availableImages : [...images];
  const img = options.splice(Math.floor(Math.random() * options.length), 1)[0];
  availableImages = options;
  return img;
}
