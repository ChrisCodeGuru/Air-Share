import React from 'react';
import ReactDOM from 'react-dom/client';
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import './css/index.css';
import Index from './pages/index';
import Signup from  './pages/signup'
import FileShare from './pages/fileshare';
import Settings from './pages/settings';
import Login from './pages/login';
import reportWebVitals from './reportWebVitals';
import { CookiesProvider } from 'react-cookie';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Index />,
  },
  {
    path: "/signup",
    element: <Signup />,
  },
  {
    path: "/fileshare",
    element: <FileShare />,
  },
  {
    path: "/settings",
    element: <Settings />,
  },
  {
    path: "/login",
    element: <Login />,
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <CookiesProvider>
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>
  </CookiesProvider>
);

reportWebVitals();
