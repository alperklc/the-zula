import React, { useEffect } from 'react';
import { useIntl } from 'react-intl';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';
import { UIProvider } from './contexts/uiContext';
import { useAuth } from './contexts/authContext';
import Home from './pages/home';
import Login from './pages/login';
import Callback from './pages/authCallback';
import { ListPage } from './pages/notes/list';
import CreateNote from './pages/notes/create';
import Note from './pages/notes/detail';
import EditNote from './pages/notes/edit';
import Profile from './pages/settings/profile';
import CreateBookmark from './pages/bookmarks/create';
import BookmarkDetails from './pages/bookmarks/detail';
import { BookmarksListPage } from './pages/bookmarks/list';
import { ActivitiesListPage } from './pages/activityLog';
import { useToast } from './components/toast/toast-message-context';
import { toStatusTextKey } from './components/toast/status-text-mapping';
import './App.css';

function PrivateRoute({ path, element }: { path: string; element: React.ReactElement }) {
  const auth = useAuth();
  const navigate = useNavigate()

  useEffect(() => {
    if(!auth.user || auth.user?.expired) {
      navigate(`/login?redirectAfterLogin=${path}`, { replace: true })
    }
  }, [auth, path, navigate])

  return auth.initialized && auth.user && !auth.user?.expired ? element : null
}

function App() {
  const { user, setSessionId } = useAuth();
  const { show: showToast } = useToast();
  const intl = useIntl();
  
  const webSocketConnection = React.useRef<WebSocket | null>(null);

  React.useEffect(() => {
    setWebSocketConnection();
    subscribeToSocketMessage();

    return () => webSocketConnection.current?.close();
  }, []);

  const setWebSocketConnection = () => {
    if (window["WebSocket"]) {
      const socketConnection = new WebSocket(`wss://${document.location.host}/api/v1/ws/${user?.profile.sub}?token=${user?.access_token}`);
      webSocketConnection.current = socketConnection;
    }
  }

  const subscribeToSocketMessage = () => {
    if (webSocketConnection.current === null) {
        return;
    }

    webSocketConnection.current.onclose = (evt) => {
      console.log('Your Connection is closed.');
    };

    webSocketConnection.current.onmessage = (event) => {
      try {
          const socketPayload = JSON.parse(event.data);
          switch (socketPayload.eventName) {
              case 'join':
              case 'disconnect':
                  if (!socketPayload.eventPayload) {
                      return
                  }
                  console.log('joined / disconnected.', socketPayload);
                  setSessionId(socketPayload.eventPayload.sessionID)
                  break;
              case 'msg':
                 const statusTextKey = toStatusTextKey(
                    socketPayload.eventPayload.resourceType,
                    socketPayload.eventPayload?.action,
                  )
                  if (statusTextKey) {
                    const message = intl.formatMessage({ id: statusTextKey })
            
                    showToast(message, 'success');
                    }
                  console.log(socketPayload)
                  
                  break;
              default:
                  break;
          }
        } catch (error) {
            console.log(error)
            console.warn('Something went wrong while decoding the Message Payload')
        }
    };
  }

  React.useEffect(() => {
    window.addEventListener('appinstalled', () => {
      alert('👍 app successfully installed')
    })

    if ('serviceWorker' in navigator) {
      window.addEventListener('load', function () {
        navigator.serviceWorker.register('/sw.js').then(
          function (registration) {
            console.log('Service Worker registration successful with scope: ', registration.scope)
          },
          function (err) {
            console.log('Service Worker registration failed: ', err)
          },
        )
      })
    }
  }, [])

  // TODO: fix theme selection
  return (
    <UIProvider initialTheme={''}>
      <Router>
        <Routes>
          <Route path="/" element={<PrivateRoute path={"/"} element={<Home />}/>} />
          <Route path="/login" element={<Login />} />
          <Route path="/callback" element={<Callback />} />
          <Route path="/activity-log" element={<PrivateRoute path={"/activity-log"} element={<ActivitiesListPage />}/>}/>
          <Route path="/notes" element={<PrivateRoute path={"/notes"} element={<ListPage />}/>}/>
          <Route path="/notes/create" element={<PrivateRoute path={"/notes/create"} element={<CreateNote />}/>}/>
          <Route path="/notes/:shortId" element={<PrivateRoute path={"/notes/:shortId"} element={<Note />}/>}/>
          <Route path="/notes/:shortId/edit" element={<PrivateRoute path={"/notes/:shortId/edit"} element={<EditNote />}/>}/>
          <Route path="/bookmarks" element={<PrivateRoute path={"/bookmarks"} element={<BookmarksListPage />}/>}/>
          <Route path="/bookmarks/create" element={<PrivateRoute path={"/bookmarks/create"} element={<CreateBookmark />}/>}/>
          <Route path="/bookmarks/:shortId" element={<PrivateRoute path={"/bookmarks/:shortId"} element={<BookmarkDetails />}/>}/>
          <Route path="/settings/profile" element={<PrivateRoute path={"/settings/profile"} element={<Profile />}/>}/>
        </Routes>
      </Router>
      </UIProvider>
    );
}

export default App
