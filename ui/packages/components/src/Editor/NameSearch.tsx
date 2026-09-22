import { Button, toast } from "@gcsim/primitives";
import type {
	IAction,
	IArtifact,
	ICharacter,
	IEnemy,
	IStat,
	IWeapon,
} from "@gcsim/types";
import { BarChart3, Bug, Footprints, Gem, Sword, Users } from "lucide-react";
import { useState } from "react";
import { Trans, useTranslation } from "react-i18next";
import {
	ActionSelect,
	ArtifactSelect,
	CharacterSelect,
	EnemySelect,
	StatSelect,
	WeaponSelect,
} from "../common/gcsim";

export function NameSearch() {
	const { t } = useTranslation();
	const [charactersOpen, setCharactersOpen] = useState(false);
	const [artifactsOpen, setArtifactsOpen] = useState(false);
	const [weaponsOpen, setWeaponsOpen] = useState(false);
	const [actionsOpen, setActionsOpen] = useState(false);
	const [statsOpen, setStatsOpen] = useState(false);
	const [enemiesOpen, setEnemiesOpen] = useState(false);

	const copyToClipboard = (item: string) => {
		navigator.clipboard.writeText(item ?? "").then(() => {
			toast.success(t("simple.copied_to_clipboard", { item }), {
				duration: 2000,
			});
		});
	};

	return (
		<div className="flex flex-col gap-1.5">
			<div className="flex flex-row gap-1.5 my-1 mx-2">
				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setCharactersOpen(true)}
				>
					<Users />
					<Trans>db.characters</Trans>
				</Button>
				<CharacterSelect
					isOpen={charactersOpen}
					onClose={() => setCharactersOpen(false)}
					onSelect={(character: ICharacter) => {
						setCharactersOpen(false);
						copyToClipboard(character);
					}}
				/>

				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setWeaponsOpen(true)}
				>
					<Sword />
					<Trans>simple.weapons</Trans>
				</Button>
				<WeaponSelect
					isOpen={weaponsOpen}
					onClose={() => setWeaponsOpen(false)}
					onSelect={(weapon: IWeapon) => {
						setWeaponsOpen(false);
						copyToClipboard(weapon);
					}}
				/>

				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setArtifactsOpen(true)}
				>
					<Gem />
					<Trans>simple.artifacts</Trans>
				</Button>
				<ArtifactSelect
					isOpen={artifactsOpen}
					onClose={() => setArtifactsOpen(false)}
					onSelect={(artifact: IArtifact) => {
						setArtifactsOpen(false);
						copyToClipboard(artifact);
					}}
				/>
			</div>
			<div className="flex flex-row gap-1.5 my-1 mx-2">
				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setEnemiesOpen(true)}
				>
					<Bug />
					<Trans>simple.enemies</Trans>
				</Button>
				<EnemySelect
					isOpen={enemiesOpen}
					onClose={() => setEnemiesOpen(false)}
					onSelect={(enemy: IEnemy) => {
						setEnemiesOpen(false);
						copyToClipboard(enemy);
					}}
				/>
				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setActionsOpen(true)}
				>
					<Footprints />
					<Trans>simple.actions</Trans>
				</Button>
				<ActionSelect
					isOpen={actionsOpen}
					onClose={() => setActionsOpen(false)}
					onSelect={(action: IAction) => {
						setActionsOpen(false);
						copyToClipboard(action);
					}}
				/>

				<Button
					variant="secondary"
					className="flex-1"
					onClick={() => setStatsOpen(true)}
				>
					<BarChart3 />
					<Trans>simple.stats</Trans>
				</Button>
				<StatSelect
					isOpen={statsOpen}
					onClose={() => setStatsOpen(false)}
					onSelect={(stat: IStat) => {
						setStatsOpen(false);
						copyToClipboard(stat);
					}}
				/>
			</div>
		</div>
	);
}
