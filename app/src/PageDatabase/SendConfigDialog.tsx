import React from 'react';
import { Dialog, Classes, Callout, Checkbox, Button } from '@blueprintjs/core';
import { Trans } from 'react-i18next';

type SendConfigDialogProps = {
  isOpen: boolean;
  handleClose: () => void;
  handleSubmit: (keep: boolean) => void;
};

const LOCALSTORAGE_KEY = 'gcsim-viewer-cpy-cfg-settings';

export function SendConfigDialog(props: SendConfigDialogProps) {
  const [keepExistingTeam, setKeepExistingTeam] = React.useState<boolean>(
    () => {
      const saved = localStorage.getItem(LOCALSTORAGE_KEY);
      if (saved === 'true') {
        return true;
      }
      return false;
    }
  );
  const handleToggleSelected = () => {
    localStorage.setItem(LOCALSTORAGE_KEY, keepExistingTeam ? 'false' : 'true');
    setKeepExistingTeam(!keepExistingTeam);
  };

  return (
    <Dialog isOpen={props.isOpen} onClose={props.handleClose}>
      <div className={Classes.DIALOG_BODY}>
        <Trans>viewer.load_this_configuration</Trans>
        <Callout intent="warning" className="mt-2">
          <Trans>viewer.this_will_overwrite</Trans>
        </Callout>
        <Checkbox
          label="Copy action list only (ignore character stats)"
          className="mt-2"
          checked={keepExistingTeam}
          onClick={handleToggleSelected}
        />
      </div>

      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button
            onClick={() => props.handleSubmit(keepExistingTeam)}
            intent="primary"
          >
            <Trans>db.continue</Trans>
          </Button>
          <Button onClick={props.handleClose}>
            <Trans>db.cancel</Trans>
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
