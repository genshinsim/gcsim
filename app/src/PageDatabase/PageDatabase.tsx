import { Spinner } from '@blueprintjs/core';
import React from 'react';
import { Viewport } from '~src/Components';
import { useAppDispatch, useAppSelector } from '~src/store';
import { DatabaseCharacters } from './DatabaseCharacters';
import { loadDB } from './dbSlice';

export function Database() {
  const status = useAppSelector((state) => state.db.status);
  const errorMsg = useAppSelector((state) => state.db.errorMsg);
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    if (status === 'idle') {
      dispatch(loadDB());
    }
  }, [status, dispatch]);

  switch (status) {
    case 'loading':
    case 'idle':
      return (
        <Viewport>
          <div className="flex flex-row place-content-center mt-2">
            <Spinner />
          </div>
        </Viewport>
      );
    case 'error':
      return (
        <Viewport>
          <div className="flex flex-row place-content-center mt-2">
            {errorMsg}
          </div>
        </Viewport>
      );
  }
  return (
    <Viewport>
      <DatabaseCharacters />
    </Viewport>
  );
}
