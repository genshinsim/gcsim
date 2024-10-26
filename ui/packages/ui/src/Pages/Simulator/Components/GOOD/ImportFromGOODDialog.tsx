import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  Position,
  Toaster,
} from '@blueprintjs/core';
import React from 'react';
import {Trans, useTranslation} from 'react-i18next';
import {useAppDispatch} from '../../../../Stores/store';
import {userDataActions} from '../../../../Stores/userDataSlice';
import {IGOODImport, parseFromGOOD} from './parseFromGOOD';

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

const lsKey = 'GOOD-import';

export function ImportFromGOODDialog(props: Props) {
  const [data, setData] = React.useState<IGOODImport>();
  const dispatch = useAppDispatch();
  const {t} = useTranslation();

  const handleLoad = () => {
    if (data !== undefined) {
      dispatch(
        userDataActions.loadFromGOOD({data: data.characters, source: 'GOOD'}),
      );
      props.onClose();
      AppToaster.show({
        message: t<string>('importer.import_success'),
        intent: 'success',
      });
    }
  };
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    localStorage.setItem(lsKey, e.target.value);
    setData(parseFromGOOD(e.target.value));
  };
  return (
    <Dialog
      className="w-screen"
      isOpen={props.isOpen}
      onClose={props.onClose}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title={t<string>('simple.tools_import', {src: 'Genshin Optimizer/GOOD'})}
      style={{width: '85%'}}>
      <div className={Classes.DIALOG_BODY}>
        <p className="!pb-2">
          <Trans i18nKey="simple.tools_import_pre_go">
            <a
              href="https://frzyc.github.io/genshin-optimizer/#/setting"
              target="_blank"
              rel="noreferrer"
            />
          </Trans>
        </p>
        <Callout intent="warning">
          {t<string>('simple.tools_import_warning', {src: 'GOOD/Enka'})}
        </Callout>
        <textarea
          value={localStorage.getItem(lsKey) ?? ''}
          onChange={handleChange}
          className="w-full p-2 bg-gray-600 rounded-md mt-2"
          rows={7}
        />
        <p className="font-bold !pt-2">
          {t<string>('simple.tools_import_after')}
        </p>
        {data ? (
          data.err === '' ? (
            <Callout intent="success" className="mt-2 p-2">
              {t<string>('simple.tools_import_post_go')}
            </Callout>
          ) : (
            <Callout intent="warning" className="mt-2 p-2">
              {data!.err}
            </Callout>
          )
        ) : null}
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <ButtonGroup>
            <Button
              onClick={handleLoad}
              disabled={!data || data.err !== ''}
              intent="primary">
              {t<string>('simple.import')}
            </Button>
            <Button onClick={props.onClose} intent="danger">
              {t<string>('db.cancel')}
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </Dialog>
  );
}
