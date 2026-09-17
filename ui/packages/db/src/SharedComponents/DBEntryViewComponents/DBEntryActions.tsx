import { dynamicKey } from "@gcsim/localization";
import { Button, toast } from "@gcsim/primitives";
import axios from "axios";
import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { AuthContext } from "../Management.context";

export default function DBEntryActions({
	id,
	share_key,
}: {
	id: string | undefined | null;
	share_key: string | undefined | null;
}) {
	const { t: translate } = useTranslation();

	const t = (key: string) => translate(dynamicKey(key));

	const isAdmin = useContext(AuthContext).isAdmin;

	if (isAdmin) {
		return (
			<div className="flex flex-col justify-center gap-2">
				<ApproveDBEntryButton dbEntryId={id} />

				<Button asChild>
					<a
						href={`https://gcsim.app/v3/viewer/share/${share_key}`}
						target="_blank"
						rel="noreferrer"
					>
						<div className="md:flex hidden">{t("db.openInViewer")}</div>
						<div className="flex md:hidden">
							{
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									strokeWidth={1.5}
									stroke="currentColor"
									className="w-5 h-5"
									role="img"
									aria-label={t("db.openInViewer") as string}
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"
									/>
								</svg>
							}
						</div>
					</a>
				</Button>
			</div>
		);
	}
	return (
		<Button asChild className="w-full">
			<a
				href={`https://gcsim.app/v3/viewer/share/${share_key}`}
				target="_blank"
				rel="noreferrer"
			>
				<div className="m-0">{t("db.openInViewer")}</div>
			</a>
		</Button>
	);
}

function ApproveDBEntryButton({
	dbEntryId,
}: {
	dbEntryId: string | undefined | null;
}) {
	return (
		<Button
			className="w-full bg-emerald-600 text-white hover:bg-emerald-600/90"
			disabled={dbEntryId === undefined || dbEntryId === null}
			onClick={() => {
				axios.post(`/api/approve/${dbEntryId}`).then((res) => {
					if (res.status === 200) {
						toast.success("Approved");
					}
				});
			}}
		>
			<div className="md:flex hidden">{"Approve"}</div>
			<div className="md:hidden flex">
				{
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						strokeWidth={1.5}
						stroke="currentColor"
						className="w-6 h-6"
						role="img"
						aria-label="Approve"
					>
						<path
							strokeLinecap="round"
							strokeLinejoin="round"
							d="M4.5 12.75l6 6 9-13.5"
						/>
					</svg>
				}
			</div>
		</Button>
	);
}
