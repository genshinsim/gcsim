import { HotkeysProvider } from '@blueprintjs/core';
import { UIReduxStore } from '@gcsim/ui';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import App from './App';

// all the css styling we need (except tailwind)
import '@blueprintjs/core/lib/css/blueprint.css';
import '@blueprintjs/icons/lib/css/blueprint-icons.css';
import './index.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Provider store={UIReduxStore}>
      <HotkeysProvider>
        <App />
      </HotkeysProvider>
    </Provider>
  </React.StrictMode>
);