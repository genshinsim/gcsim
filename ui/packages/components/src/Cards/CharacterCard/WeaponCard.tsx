import type { model } from "@gcsim/types";

type Props = {
	weapon: model.Weapon;
	name: string;
	isSkeleton?: boolean;
};

export function WeaponCard({ weapon, name, isSkeleton }: Props) {
	const content = (
		<div className="flex flex-row">
			<div className="w-12 h-12">
				<img
					src={`/api/assets/weapons/${weapon.name}.png`}
					alt={weapon.name}
					className="object-contain w-full"
					onError={(e) =>
						((e.target as HTMLImageElement).src =
							"/api/assets/misc/default.png")
					}
				/>
			</div>
			<div className="flex-grow text-g-sm pl-2 flex flex-col justify-center">
				<div className="font-medium text-left">
					{name.replace(/(.{20})..+/, "$1…") + " R" + weapon.refine}
				</div>
				<div className="justify-center items-center rounded-g-md">
					Lvl {weapon.level}/{weapon.max_level}
				</div>
			</div>
		</div>
	);

	return (
		<div className="weapon-parent ml-2 mr-2 p-2 bg-g-surface-2 border-g-line border">
			{isSkeleton ? <div className="h-12"></div> : content}
		</div>
	);
}
