import {AnchorButton, Card} from '@blueprintjs/core';
import {useTranslation} from 'react-i18next';
import {useEffect, useState} from 'react';
import axios from 'axios';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { AvatarPortrait } from "@gcsim/components";
import LatestCharactersData from '@gcsim/data/src/latest_chars.json'

const majorVersionRegex = /v\d+\.\d+/gm;

export function LatestVersion() {
  const {t} = useTranslation();

  const [{isLoaded, text, tag, portraits}, setState] = useState({
    isLoaded: false,
    text: '',
    tag: '',
    portraits: []
  });

  useEffect(() => {
    axios('https://api.github.com/repos/genshinsim/gcsim/releases/latest')
      .then((resp: {data}) => {
        const majorVersion = majorVersionRegex.exec(resp.data.name);
        let portraits = [];
        if (majorVersion && majorVersion[0]) {
          portraits = LatestCharactersData[majorVersion[0]] || [];
        }
        setState({isLoaded: true, text: resp.data.body, tag: resp.data.name, portraits});
      })
      .catch((err) =>
        console.log(t('viewer.error_encountered') + err.message),
      );
  }, [t]);

  return (
    <Card className="flex flex-col items-center gap-4 overflow-x-auto">
      {isLoaded ? (
        <>
          <div className="flex flex-col gap-4">
            <h1 className="text-center text-xl md:text-2xl lg:text-4xl">
              <b>{t('dash.latest_release')}</b>
              <a
                href={`https://github.com/genshinsim/gcsim/releases/tag/${tag}`}>
                {tag}
              </a>
            </h1>
          </div>
          <div className="flex flex-col">
            <h2 className="text-center text-2xl">{t('dash.new_characters')}</h2>
            <div className="flex gap-4">
              {portraits.map((char) => (
                <div key={char} className="flex flex-col items-center">
                  <AvatarPortrait char={{ name: char }} hideDetails />
                  {t(`game:character_names.${char}`)}
                </div>
              ))}
            </div>
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
            className="!p-3 !rounded-md">
            <span className="text-xl md:text-2xl font-semibold">
              {t('dash.view_releases')}
            </span>
          </AnchorButton>
        </>
      ) : (
        <>{t('sim.loading')}</>
      )}
    </Card>
  );
}