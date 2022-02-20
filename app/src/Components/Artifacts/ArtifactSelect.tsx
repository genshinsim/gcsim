import { Omnibar } from "@blueprintjs/select";
import { IArtifact, artifactSelectProps } from "./artifacts";

const ArtifactOmnibar = Omnibar.ofType<IArtifact>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: IArtifact) => void;
};
export function ArtifactSelect(props: Props) {
  return (
    <ArtifactOmnibar
      resetOnSelect
      {...artifactSelectProps}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
