import { Omnibar } from "@blueprintjs/select";
import { IArtifact } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { artifactSelectProps } from "./artifacts";

const ArtifactOmnibar = Omnibar.ofType<IArtifact>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (artifact: IArtifact) => void;
};

export function ArtifactSelect(props: Props) {
  const { t } = useTranslation();
  return (
    <ArtifactOmnibar
      resetOnSelect
      {...artifactSelectProps}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
      className="!absolute !left-0 !right-0 !mx-auto !w-80"
    />
  );
}
