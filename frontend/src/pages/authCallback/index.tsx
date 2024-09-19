import { useEffect } from "react";
import { User, UserManager } from "oidc-client-ts";
import { useAuth } from "../../contexts/authContext";
import { useNavigate } from 'react-router-dom'
import LoadingIndicator from "../../components/loadingIndicator";

async function redirectAfterLogin(user: User | null, userManager?: UserManager) {
  if (!userManager) {
    return
  }

  let userData: User | null = null
  try {
    if (user === null) {
      userData = await userManager.signinRedirectCallback()
    } else {
      userData = await userManager.getUser()
    }
  } catch (error: unknown) {
    console.error(error)
    userData = null;

  }
  return userData
}

const Callback = () => {
  const navigate = useNavigate()
  const { user, setUser, userManager } = useAuth()

  useEffect(() => {
    redirectAfterLogin(user, userManager)
      .then(user => { 
        setUser(user ?? null)
      })
      .catch(err => {
        console.error(err)
        setUser(null)
      })
  }, []);

  useEffect(() => {
    if (user?.url_state) {
      navigate(user?.url_state, { replace: true })
    }
  }, [user?.url_state, navigate])

  return <div><LoadingIndicator /></div>
};

export default Callback;