import React from 'react'
import { User, UserProfile } from "oidc-client-ts";
import { RawIntlProvider, createIntl, createIntlCache } from 'react-intl'
import { useAuth } from './authContext'
import English from "../messages/en.json"
import German from "../messages/de.json"

interface UIContext {
  language: string
  theme: string
  isMobile: boolean
  backdropActive: boolean
  switchLanguage: (_: string) => void
  switchTheme: (_: string) => void
  toggleBackdrop: () => void
}

export const UIContext = React.createContext<UIContext>({
  language: 'en',
  theme: 'darkTheme',
  isMobile: false,
  backdropActive: false,
  switchLanguage: (_: string) => ({}),
  switchTheme: (_: string) => ({}),
  toggleBackdrop: () => ({}),
})

const messages = {
  'en': English,
  'de': German
}

const initialLocale = 'en'

export const cache = createIntlCache()

export let intl = createIntl(
  {locale: initialLocale, messages: English},
  cache
)

type ExtendedProfile = UserProfile & {
  "urn:zitadel:iam:user:metadata"?: {theme?: string}[]
}

function getTheme(user: User | null) {
  const metadata = (user?.profile as ExtendedProfile)?.["urn:zitadel:iam:user:metadata"]
  if (!metadata || metadata.length === 0) {
    return "light"
  }

  const theme = metadata[0].theme ?? ""
  return atob(theme)
}

export const UIProvider = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuth()
  const defaultLocale = user?.profile.locale || 'en'
  // TODO: read initial theme here
 
  const isMobileDevice = !!(navigator.userAgent || '').match(
    /Android|BlackBerry|iPhone|iPad|iPod|Opera Mini|IEMobile|WPDesktop/i,
  )

  const [language, setLanguage] = React.useState<string>(defaultLocale)
  const [theme, setTheme] = React.useState<string>(getTheme(user!))
  const [backdropActive, setBackdropActive] = React.useState<boolean>(false)
  const [isMobile, setIsMobile] = React.useState<boolean>(isMobileDevice)

  const switchLanguage = (nextLanguage: string) => {
    document.documentElement.setAttribute('lang', nextLanguage)

    intl = createIntl({locale: nextLanguage, messages: messages[nextLanguage as 'en' |'de']}, cache)
    document.documentElement.lang = nextLanguage

    setLanguage(nextLanguage)
  }

  const switchTheme = (nextTheme: string) => {
    document.body.setAttribute('data-theme', nextTheme)
    setTheme(nextTheme)
  }

  const toggleBackdrop = () => {
    setBackdropActive(!backdropActive)
  }
  
  const handleResize = () => {
    setIsMobile(window.innerWidth < 700)
  }

  React.useEffect(() => {
    if (!isMobileDevice) {
      window.addEventListener("resize", handleResize);
      
      return () => {
        window.removeEventListener("resize", handleResize)
      }
    }
  }, [isMobileDevice]);

  return (
    <UIContext.Provider
      value={{
        language,
        theme,
        backdropActive,
        isMobile,
        switchLanguage,
        switchTheme,
        toggleBackdrop,
      }}
    >
      <RawIntlProvider value={intl}>
        {children}
      </RawIntlProvider>
    </UIContext.Provider>
  )
}

export function useUI() {
  const context = React.useContext(UIContext)
  if (context === undefined) {
    throw new Error('useUI must be used within an UIProvider')
  }
  return context
}
