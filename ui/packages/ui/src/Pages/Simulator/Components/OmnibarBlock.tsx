import { Button, Intent, Position, Toaster } from "@blueprintjs/core";
import { IAction, IArtifact, ICharacter, IStat, IWeapon } from "@gcsim/types";
import { ArtifactSelect, WeaponSelect } from "@ui/Components/Select";
import { ActionSelect } from "@ui/Components/Select/ActionSelect";
import { CharacterSelect } from "@ui/Components/Select/CharacterSelect";
import { StatSelect } from "@ui/Components/Select/StatSelect";
import { useRef, useState } from "react";
import { Trans, useTranslation } from "react-i18next";

export function OmnibarBlock() {
  const { t } = useTranslation();
  const [charactersOpen, setCharactersOpen] = useState(false);
  const [artifactsOpen, setArtifactsOpen] = useState(false);
  const [weaponsOpen, setWeaponsOpen] = useState(false);
  const [actionsOpen, setActionsOpen] = useState(false);
  const [statsOpen, setStatsOpen] = useState(false);

  const copyToast = useRef<Toaster>(null);

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex flex-row gap-1.5 my-1 mx-2">
        <Button
          icon="people"
          fill
          onClick={() => {
            setCharactersOpen(true);
          }}
        >
          <Trans>db.characters</Trans>
        </Button>
        <CharacterSelect
          isOpen={charactersOpen}
          onClose={() => setCharactersOpen(false)}
          onSelect={(character: ICharacter) => {
            setCharactersOpen(false);
            navigator.clipboard.writeText(character ?? "").then(() => {
              copyToast.current?.show({
                message: `${t("simple.copied_to_clipboard", {
                  item: character,
                })}`,
                intent: Intent.SUCCESS,
                timeout: 2000,
              });
            });
          }}
        />

        <Button
          icon="build"
          fill
          onClick={() => {
            setWeaponsOpen(true);
          }}
        >
          <Trans>simple.weapons</Trans>
        </Button>
        <WeaponSelect
          isOpen={weaponsOpen}
          onClose={() => setWeaponsOpen(false)}
          onSelect={(weapon: IWeapon) => {
            setWeaponsOpen(false);
            navigator.clipboard.writeText(weapon ?? "").then(() => {
              copyToast.current?.show({
                message: `${t("simple.copied_to_clipboard", { item: weapon })}`,
                intent: Intent.SUCCESS,
                timeout: 2000,
              });
            });
          }}
        />

        <Button
          icon="glass"
          fill
          onClick={() => {
            setArtifactsOpen(true);
          }}
        >
          <Trans>simple.artifacts</Trans>
        </Button>
        <ArtifactSelect
          isOpen={artifactsOpen}
          onClose={() => setArtifactsOpen(false)}
          onSelect={(artifact: IArtifact) => {
            setArtifactsOpen(false);
            navigator.clipboard.writeText(artifact ?? "").then(() => {
              copyToast.current?.show({
                message: `${t("simple.copied_to_clipboard", {
                  item: artifact,
                })}`,
                intent: Intent.SUCCESS,
                timeout: 2000,
              });
            });
          }}
        />
      </div>
      <div className="flex flex-row gap-1.5 my-1 mx-2">
        <Button
          icon="walk"
          fill
          onClick={() => {
            setActionsOpen(true);
          }}
        >
          <Trans>simple.actions</Trans>
        </Button>
        <ActionSelect
          isOpen={actionsOpen}
          onClose={() => setActionsOpen(false)}
          onSelect={(action: IAction) => {
            setActionsOpen(false);
            navigator.clipboard.writeText(action ?? "").then(() => {
              copyToast.current?.show({
                message: `${t("simple.copied_to_clipboard", { item: action })}`,
                intent: Intent.SUCCESS,
                timeout: 2000,
              });
            });
          }}
        />

        <Button
          icon="panel-stats"
          fill
          onClick={() => {
            setStatsOpen(true);
          }}
        >
          <Trans>simple.stats</Trans>
        </Button>
        <StatSelect
          isOpen={statsOpen}
          onClose={() => setStatsOpen(false)}
          onSelect={(stat: IStat) => {
            setStatsOpen(false);
            navigator.clipboard.writeText(stat ?? "").then(() => {
              copyToast.current?.show({
                message: `${t("simple.copied_to_clipboard", { item: stat })}`,
                intent: Intent.SUCCESS,
                timeout: 2000,
              });
            });
          }}
        />
      </div>
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
    </div>
  );
}
