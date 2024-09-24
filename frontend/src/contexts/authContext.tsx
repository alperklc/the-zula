/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react'
import { createZitadelAuth, ZitadelAuth } from '@zitadel/react'
import { User, UserManager, UserProfile } from "oidc-client-ts";
import { useTranslation } from 'react-i18next';
import { FrontendConfig } from '../types/Api';

interface State {
  user: User | null
  sessionId: string,
  setSessionId: (_: string) => void,
  initialized: boolean,
  setUser: (_: User | null) => void,
  login: (_: string) => void,
  logout: () => void,
  userManager?: UserManager,
}

export const AuthContext = React.createContext<State>({
  user: null,
  sessionId: "",
  setSessionId: () => ({}),
  initialized: false,
  setUser: (_: User | null) => ({}),
  login: () => ({}),
  logout: () => ({}),
})

export const useAuth = () => React.useContext(AuthContext)

type ExtendedProfile = UserProfile & {
  "urn:zitadel:iam:user:metadata"?: {theme?: string}[]
}

function getTheme(user: User | null) {
  const metadata = (user?.profile as ExtendedProfile)?.["urn:zitadel:iam:user:metadata"]
  if (!metadata || metadata.length === 0) {
    return "lightTheme"
  }

  const theme = metadata[0].theme ?? ""
  return atob(theme)
}

export const AuthContextProvider = ({ children, initialConfig }: { children: JSX.Element, initialConfig: FrontendConfig }) => {
  const { i18n } = useTranslation()

  const [initialized, setInitialized] = React.useState<boolean>(false)
  const [sessionId, setSessionId] = React.useState<string>("")
  const [user, setuser] = React.useState<User | null>(null)
  const zitadel = React.useRef<ZitadelAuth>()
  
  React.useEffect(() => {
    zitadel.current = createZitadelAuth(initialConfig);
  }, [initialConfig]);
  
  function login(redirectTo: string) {
    zitadel.current?.userManager.signinRedirect({redirect_uri: initialConfig.redirect_uri, url_state: redirectTo, scope: "openid profile email urn:zitadel:iam:user:metadata"});
    zitadel.current?.userManager.startSilentRenew()
  }

  function logout() {
    zitadel.current?.signout();
  }

  const setUserProfile = (user: User | null) => {
    setuser(!user?.expired ? user : null);
    i18n.changeLanguage(user?.profile.locale)
  }

  const setUser = (user: User | null) => {
    setuser(!user?.expired ? user : null);
    i18n.changeLanguage(user?.profile.locale)
    const theme = getTheme(user)
    document.body.setAttribute('data-theme', theme)
  }

  React.useEffect(() => {
    zitadel.current?.userManager.getUser()
      .then(setUserProfile)
      .finally(() => {
      setInitialized(true)
    });
  }, [zitadel.current]);
 
  return (
    <AuthContext.Provider
      value={{
        user,
        sessionId,
        setSessionId,
        initialized,
        setUser,
        login,
        logout,
        userManager: zitadel.current?.userManager,
      }}
    >
      {initialized ? children : <></>}
    </AuthContext.Provider>
  )
}
