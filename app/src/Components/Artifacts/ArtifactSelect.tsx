import { Omnibar } from "@blueprintjs/select";
import { artifactSelectProps, IArtifact } from "./artifacts";
import { useTranslation } from 'react-i18next'

const ArtifactOmnibar = Omnibar.ofType<IArtifact>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: IArtifact) => void;
};

export function ArtifactSelect(props: Props) {
  let { i18n } = useTranslation()

  return (
    <ArtifactOmnibar
      resetOnSelect
      {...artifactSelectProps[i18n.language]}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
