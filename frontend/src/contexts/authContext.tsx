/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react'
import { ZitadelConfig, createZitadelAuth } from '@zitadel/react'
import { User, UserManager } from "oidc-client-ts";

interface State {
  user: User | null
  initialized: boolean,
  setUser: (_: User | null) => void,
  login: (_: string) => void,
  logout: () => void,
  userManager?: UserManager,
}

export const AuthContext = React.createContext<State>({
  user: null,
  initialized: false,
  setUser: (_: User | null) => ({}),
  login: () => ({}),
  logout: () => ({}),
})

export const useAuth = () => React.useContext(AuthContext)

const config: ZitadelConfig = {
  authority: "https://auth.local.the-zula.app:8080",
  client_id: "265970737444093958@zula",
  // client_id: "265972414259724294@zula",
  //"261045606791839750@the-zula",
  redirect_uri: "https://local.the-zula.app/callback",
  post_logout_redirect_uri: "https://local.the-zula.app",
};

export const AuthContextProvider = ({ children }: { children: JSX.Element }) => {
  const [initialized, setInitialized] = React.useState<boolean>(false)
  const [user, setUser] = React.useState<User | null>(null)
  const zitadel = createZitadelAuth(config);
  
  function login(redirectTo: string) {
    zitadel.userManager.signinRedirect({redirect_uri: config.redirect_uri, url_state: redirectTo, scope: "openid profile email urn:zitadel:iam:user:metadata"});
  }

  function logout() {
    zitadel.signout();
  }

  React.useEffect(() => {
    zitadel.userManager.getUser().then((user) => {
      setUser(!user?.expired ? user : null);
    }).finally(() => {
      setInitialized(true)
    });
  }, []);
 
  return (
    <AuthContext.Provider
      value={{
        user,
        initialized,
        setUser,
        login,
        logout,
        userManager: zitadel.userManager,
      }}
    >
      {initialized ? children : <></>}
    </AuthContext.Provider>
  )
}
