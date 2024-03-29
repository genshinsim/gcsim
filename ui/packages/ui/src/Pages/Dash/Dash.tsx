import { AnchorButton, Card } from "@blueprintjs/core";
import DBEntryView from "@gcsim/db/src/SharedComponents/DBEntryView";
import { db } from "@gcsim/types";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router-dom";
import remarkGfm from "remark-gfm";

const randQuery = {
  query: {
    $sampleRate: 0.01,
  },
  limit: 1,
  skip: 0,
  sort: {
    create_date: -1,
  },
};

export function Dash() {
  const { t } = useTranslation();

  const [{ isLoaded, text, tag }, setState] = useState({
    isLoaded: false,
    text: "",
    tag: "",
  });
  const [data, setData] = useState<db.Entry[]>([]);
  const [dataIsLoaded, setDataIsLoaded] = useState(false);

  useEffect(() => {
    axios("https://api.github.com/repos/genshinsim/gcsim/releases/latest")
      .then((resp: { data }) => {
        setState({ isLoaded: true, text: resp.data.body, tag: resp.data.name });
      })
      .catch((err) =>
        console.log(t<string>("viewer.error_encountered") + err.message)
      );
    axios(`/api/db?q=${encodeURIComponent(JSON.stringify(randQuery))}`)
      .then((resp: { data: db.Entries }) => {
        if (resp.data && resp.data.data) {
          setData(resp.data.data);
          setDataIsLoaded(true);
        }
      })
      .catch((err) =>
        console.log(t<string>("viewer.error_encountered") + err.message)
      );
  }, [t]);

  return (
    <main className="w-full flex flex-col items-center flex-grow gap-4 py-4 px-4">
      <Link
        to="/simulator"
        role="button"
        className="bp4-button bp4-intent-success !p-3 !rounded-md"
        tabIndex={0}
      >
        <span className="bp4-button-text text-3xl md:text-4xl lg:text-5xl font-semibold">
          {t<string>("dash.get_started")}
        </span>
      </Link>
      {/* mobile Chrome/Safari needs w-full for padding to work properly, mobile Firefox works fine though... */}
      <div className="flex flex-col gap-4 w-full md:w-fit">
        <Card className="flex flex-col gap-4 items-center">
          <h1 className="text-center text-xl md:text-2xl lg:text-4xl">
            <b>{t<string>("dash.users_submitted")}</b>
          </h1>
          <div>
            {dataIsLoaded
              ? data.map((entry, index) => (
                  <DBEntryView dbEntry={entry} key={index} />
                ))
              : t<string>("sim.loading")}
          </div>
          <AnchorButton
            href="https://simpact.app/"
            intent="primary"
            target="_blank"
            className="!p-3 !rounded-md"
          >
            <span className="text-xl md:text-2xl font-semibold">
              {t<string>("dash.visit_teams_db")}
            </span>
          </AnchorButton>
        </Card>
        <Card className="flex flex-col items-center gap-4 overflow-x-auto">
          {isLoaded ? (
            <>
              <div className="flex flex-col gap-4">
                <h1 className="text-center text-xl md:text-2xl lg:text-4xl">
                  <b>{t<string>("dash.latest_release")}</b>
                  <a
                    href={`https://github.com/genshinsim/gcsim/releases/tag/${tag}`}
                  >
                    {tag}
                  </a>
                </h1>
              </div>
              <div className="self-start">
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {text}
                </ReactMarkdown>
              </div>
              <AnchorButton
                href="https://github.com/genshinsim/gcsim/releases"
                intent="primary"
                target="_blank"
                className="!p-3 !rounded-md"
              >
                <span className="text-xl md:text-2xl font-semibold">
                  {t<string>("dash.view_releases")}
                </span>
              </AnchorButton>
            </>
          ) : (
            t<string>("sim.loading")
          )}
        </Card>
      </div>
    </main>
  );
}
