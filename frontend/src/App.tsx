import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';
import { UIProvider } from './contexts/uiContext';
import { useAuth } from './contexts/authContext';
import Home from './pages/home';
import Login from './pages/login';
import Callback from './pages/authCallback';
import { ListPage } from './pages/notes/list';
import CreateNote from './pages/notes/create';
import Note from './pages/notes/detail';
import './App.css'
import Profile from './pages/settings/profile';
import CreateBookmark from './pages/bookmarks/create';
import BookmarkDetails from './pages/bookmarks/detail';
import { BookmarksListPage } from './pages/bookmarks/list';

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
  return (
    <UIProvider initialTheme={''}>
      <Router>
        <Routes>
          <Route path="/" element={<PrivateRoute path={"/"} element={<Home />}/>} />
          <Route path="/login" element={<Login />} />
          <Route path="/callback" element={<Callback />} />
          <Route path="/notes" element={<PrivateRoute path={"/notes"} element={<ListPage />}/>}/>
          <Route path="/notes/create" element={<PrivateRoute path={"/notes/create"} element={<CreateNote />}/>}/>
          <Route path="/notes/:noteId" element={<PrivateRoute path={"/notes/:noteId"} element={<Note />}/>}/>
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
