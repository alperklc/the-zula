import ReactDOM from 'react-dom/client'
import { AuthContextProvider } from './contexts/authContext.tsx'
import App from './App.tsx'
import { ToastMessageProvider } from './components/toast/toast-message-context.tsx'
import { createIntl, createIntlCache, RawIntlProvider } from 'react-intl'
import English from "./messages/en.json"
import './index.css'

// TODO: fix language switching
export const cache = createIntlCache()
export const intl = createIntl(
  {locale: "en", messages: English},
  cache
)

ReactDOM.createRoot(document.getElementById('root')!).render(
  <>
    <AuthContextProvider>
      <RawIntlProvider value={intl}>
        <ToastMessageProvider>
          <App />
        </ToastMessageProvider>
      </RawIntlProvider>
    </AuthContextProvider>
  </>,
)
