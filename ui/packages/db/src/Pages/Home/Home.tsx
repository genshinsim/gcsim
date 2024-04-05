import {Button, ButtonGroup} from '@blueprintjs/core';
import tagData from '@gcsim/data/src/tags.json';
import {useTranslation} from 'react-i18next';
import {useLocation} from 'wouter';

export const Home = () => {
  const {t} = useTranslation();

  const [_, to] = useLocation();
  const sortedTagnames = Object.keys(tagData)
    .filter((key) => {
      return key !== '0' && key != '2';
    })
    .map((key) => {
      let name = tagData[key]['display_name'];
      if (key == '1') {
        name = '(Not Tagged)';
      }
      return (
        <li key={key}>
          <span className="font-semibold text-rose-700">{name}</span>
          {`: ${t<string>('db.home.tag_desc_' + key)}`}
        </li>
      );
    });
  return (
    <div className="ml-2 mr-2 mt-2">
      <div className="text-center text-lg font-semibold text-indigo-600">
        {t<string>('db.home.welcome')}{' '}
      </div>
      <div className="mb-4">
        <p className="m-2">{t<string>('db.home.simpact_desc')}</p>
        <p className="m-2">{t<string>('db.home.simpact_tag_desc')} </p>
        <p className="m-2">{t<string>('db.home.simpact_tag_list_header')} </p>
        <ul className="list-disc m-4 ml-8">{sortedTagnames}</ul>
      </div>
      <ButtonGroup fill className="mb-4">
        <Button intent="primary" onClick={() => to('/database')}>
          {t<string>('db.home.get_started')}
        </Button>
      </ButtonGroup>
    </div>
  );
};
