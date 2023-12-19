import { Button, Card, Collapse, Elevation } from "@blueprintjs/core";
import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router-dom";
import remarkGfm from "remark-gfm";

interface DashCardProps {
  children: React.ReactNode;
  href: string;
  target?: string;
}

function DashCard({ children, href, target }: DashCardProps) {
  return (
    <div className="main-page-button-container">
      {target ? (
        <a href={href} target={target}>
          <Card
            interactive
            elevation={Elevation.TWO}
            className="main-page-card"
          >
            {children}
          </Card>
        </a>
      ) : (
        <Link to={href}>
          <a>
            <Card
              interactive
              elevation={Elevation.TWO}
              className="main-page-card"
            >
              {children}
            </Card>
          </a>
        </Link>
      )}
    </div>
  );
}

export function Dash() {
  useTranslation();

  const [{ isLoaded, text, tag }, setState] = useState({
    isLoaded: false,
    text: "",
    tag: "",
  });
  const [tagIsOpen, setTagIsOpen] = useState(false);

  // for size testing use: https://api.github.com/repos/genshinsim/gcsim/releases/tags/<tag name>
  useEffect(() => {
    fetch(`https://api.github.com/repos/genshinsim/gcsim/releases/latest`)
      .then((res) => res.arrayBuffer())
      .then((buffer) => {
        const decoder = new TextDecoder("utf-8");
        const data = decoder.decode(buffer);
        const release = JSON.parse(data);
        setState({ isLoaded: true, text: release.body, tag: release.name });
      })
      .catch((err) => console.log("Error: " + err.message));
  }, []);

  return (
    <main className="w-full flex flex-col items-center flex-grow gap-4 mt-2">
      <div className="flex items-center justify-center w-full flex-grow text-2xl md:text-4xl lg:text-6xl px-4 text-center">
        <h1 className="max-w-sm md:max-w-lg lg:max-w-4xl">
          <b>
            gcsim is a Team DPS / Combat Simulation Tool for Genshin Impact.
          </b>
        </h1>
      </div>
      <div className="flex flex-col flex-grow items-center px-8 mb-4">
        {isLoaded ? (
          <>
            <div className="flex flex-col gap-4 mb-4">
              <h1 className="text-center text-xl md:text-2xl lg:text-4xl">
                <b>Latest Release: </b>
                <a
                  href={`https://github.com/genshinsim/gcsim/releases/tag/${tag}`}
                >
                  {tag}
                </a>
              </h1>
              <Button
                className="w-[100%]"
                onClick={() => setTagIsOpen(!tagIsOpen)}
              >
                {tagIsOpen ? "Hide" : "Show"} release notes
              </Button>
            </div>
            <Collapse
              className="text-md"
              isOpen={tagIsOpen}
              keepChildrenMounted={true}
            >
              <ReactMarkdown children={text} remarkPlugins={[remarkGfm]} />
            </Collapse>
          </>
        ) : (
          "Loading..."
        )}
      </div>
    </main>
  );
}
