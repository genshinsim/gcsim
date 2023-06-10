import { Position, Toaster } from "@blueprintjs/core";
import axios from "axios";
import { useContext } from "react";
import { ReactI18NextChild, useTranslation } from "react-i18next";
import { AuthContext } from "../Management.context";

export default function DBEntryActions({
  id,
  simulation_key,
}: {
  id: string | undefined | null;
  simulation_key: string | undefined | null;
}) {
  const { t: translate } = useTranslation();

  const t = (key: string) => translate(key) as ReactI18NextChild; // idk why this is needed

  const isAdmin = useContext(AuthContext).isAdmin;

  if (isAdmin) {
    return (
      <div className="flex flex-col justify-center gap-2">
        <ApproveDBEntryButton dbEntryId={id} />

        <a
          href={`https://gcsim.app/v3/viewer/share/${simulation_key}`}
          target="_blank"
          className="bp4-button    bp4-intent-primary"
          rel="noreferrer"
        >
          <div className="md:flex hidden">
            {t("db.openInViewer") as ReactI18NextChild}
          </div>
          <div className="flex md:hidden">
            {
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-5 h-5"
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
      </div>
    );
  }
  return (
    <div className="flex flex-col justify-center">
      <a
        href={`https://gcsim.app/v3/viewer/share/${simulation_key}`}
        target="_blank"
        className="bp4-button    bp4-intent-primary "
        rel="noreferrer"
      >
        <div className="md:block hidden m-0">{t("db.openInViewer")}</div>
        <div className="flex md:hidden">
          {
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-5 h-5"
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
    </div>
  );
}

function ApproveDBEntryButton({
  dbEntryId,
}: {
  dbEntryId: string | undefined | null;
}) {
  return (
    <>
      <button
        className="bp4-button bp4-intent-success w-full"
        disabled={dbEntryId === undefined || dbEntryId === null}
        onClick={() => {
          axios.post(`/api/approve/${dbEntryId}`).then((res) => {
            if (res.status === 200) {
              AppToaster.show({ message: "Approved" });
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
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M4.5 12.75l6 6 9-13.5"
              />
            </svg>
          }
        </div>
      </button>
    </>
  );
}

const AppToaster = Toaster.create({
  className: "recipe-toaster",
  position: Position.TOP,
});
