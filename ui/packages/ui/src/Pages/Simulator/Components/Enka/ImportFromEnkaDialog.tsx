import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  Position,
  Toaster,
} from '@blueprintjs/core';
import {Character} from '@gcsim/types';
import React from 'react';
import {Trans, useTranslation} from 'react-i18next';
import {useAppDispatch} from '../../../../Stores/store';
import {userDataActions} from '../../../../Stores/userDataSlice';
import FetchCharsFromEnka from './FetchCharsFromEnka';

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const AppToaster = Toaster.create({
  position: Position.BOTTOM,
});

const lsKey = 'Enka-UID';

interface CustDivProps
  extends React.DetailedHTMLProps<
    React.HTMLAttributes<HTMLDivElement>,
    HTMLDivElement
  > {
  i18nIsDynamicList?: boolean;
}

const CustDiv: React.FC<CustDivProps> = ({children, ...props}) => {
  return <div {...props}>{children}</div>;
};

export function ImportFromEnkaDialog(props: Props) {
  const {t} = useTranslation();
  const [message, setMessage] = React.useState<string>('');
  const [errors, setErrors] = React.useState<string[]>([]);
  const [characters, setCharacters] = React.useState<Character[]>([]);
  const [uid, setUid] = React.useState<string>('');
  const dispatch = useAppDispatch();

  async function handleClick() {
    localStorage.setItem(lsKey, uid);
    if (uid && validateUid(uid)) {
      try {
        setCharacters([]);
        const result = await FetchCharsFromEnka(uid);
        setErrors(result.errors ? result.errors : []);
        console.log(result);
        dispatch(
          userDataActions.loadFromGOOD({
            data: result.characters,
            source: 'enka',
          }),
        );
        setMessage('success');
        setCharacters(result.characters);
      } catch (e) {
        setMessage(`Error importing chars: ${e}`);
      }
    } else {
      setMessage('Invalid UID');
    }
  }

  return (
    <Dialog
      className="w-screen"
      isOpen={props.isOpen}
      onClose={() => {
        props.onClose();
        setMessage('');
      }}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title={t<string>('simple.tools_import', {src: 'Enka.Network'})}
      style={{width: '85%'}}>
      <div className={Classes.DIALOG_BODY}>
        <p className="!pb-2">
          <Trans i18nKey="simple.tools_import_pre_enka">
            <a href="https://enka.network/" target="_blank" rel="noreferrer" />
          </Trans>
        </p>
        <Callout intent="warning">
          {t<string>('simple.tools_import_warning', {src: 'GOOD/Enka'})}
        </Callout>
        <input
          value={uid}
          onChange={(e) => {
            setUid(e.target.value.trim());
          }}
          className="w-full p-2 bg-gray-600 rounded-md mt-2"
          placeholder={t<string>('simple.tools_paste_uid')}
        />

        {message === 'success' ? (
          <>
            <Callout intent="success" className="mt-2 p-2">
              {characters.length > 0 ? (
                <Trans i18nKey="simple.tools_import_post_enka">
                  <CustDiv i18nIsDynamicList>
                    {characters.map((e, i) => {
                      return (
                        <div key={i} className="ml-2">
                          {e.name}{' '}
                          {e.enka_build_name
                            ? '(' + e.enka_build_name + ')'
                            : ''}
                        </div>
                      );
                    })}
                  </CustDiv>
                </Trans>
              ) : null}
            </Callout>
            {errors.length > 0 ? (
              <Callout intent="warning" className="mt-2 p-2">
                Encountered the following issue(s) importing data:
                {errors.map((e, i) => {
                  return (
                    <div key={i} className="ml-2">
                      {e}
                    </div>
                  );
                })}
              </Callout>
            ) : null}
          </>
        ) : (
          <div>
            {message && (
              <Callout intent="warning" className="mt-2 p-2">
                {message}
              </Callout>
            )}
          </div>
        )}

        <p className="font-bold !pt-2">
          {t<string>('simple.tools_import_after')}
        </p>
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <ButtonGroup>
            <Button onClick={handleClick} intent="primary">
              {t<string>('simple.import')}
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </Dialog>
  );
}

function validateUid(uid: string) {
  if (!/^(18|[1-35-9])\d{8}$/.test(uid)) {
    AppToaster.show({
      message: 'Invalid UID',
      intent: 'danger',
    });
    return false;
  }
  return true;
}
