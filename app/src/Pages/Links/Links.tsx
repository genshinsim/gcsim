import { Trans, useTranslation } from 'react-i18next'

export function Links() {
  useTranslation()

  return (
    <div className="m-2 text-center text-lg">
      <Trans>links.error_loading_database</Trans>
    </div>
  );
}
