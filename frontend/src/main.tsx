import ReactDOM from 'react-dom/client'
import { AuthContextProvider } from './contexts/authContext.tsx'
import { ToastMessageProvider } from './components/toast/toast-message-context.tsx'
import { Api } from './types/Api.ts'
import App from './App.tsx'
import './i18n'
import './index.css'

const api = new Api()
const configResponse = await api.api.getFrontendConfig()
const feConfig = await configResponse.json()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <>
    <AuthContextProvider initialConfig={feConfig}>
        <ToastMessageProvider>
          <App />
        </ToastMessageProvider>
    </AuthContextProvider>
  </>,
)
