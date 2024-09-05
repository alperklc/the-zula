import ReactDOM from 'react-dom/client'
import { AuthContextProvider } from './contexts/authContext.tsx'
import { ToastMessageProvider } from './components/toast/toast-message-context.tsx'
import App from './App.tsx'
import './i18n'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <>
    <AuthContextProvider>
        <ToastMessageProvider>
          <App />
        </ToastMessageProvider>
    </AuthContextProvider>
  </>,
)
