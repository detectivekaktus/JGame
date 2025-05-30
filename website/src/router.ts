import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { LoginPage } from "./pages/LoginPage";
import { SignupPage } from "./pages/SignupPage";
import { MainPage } from "./pages/MainPage";
import { ProfilePage } from "./pages/ProfilePage";
import { SettingsPage } from "./pages/SettingsPage";

export const APP_ROUTER = createBrowserRouter([
  {
    path: "/",
    children: [
      { index: true, Component: HomePage },
      { path: "main", Component: MainPage },
      {
        path: "auth",
        children: [
          { index: true, Component: NotFoundPage },
          { path: "login", Component: LoginPage },
          { path: "signup", Component: SignupPage }
        ]
      },
      {
        path: "profile",
        children: [
          { index: true, Component: NotFoundPage },
          { path: ":id", Component: ProfilePage },
          { path: "settings", Component: SettingsPage }
        ]
      },
      { path: "*", Component: NotFoundPage }
    ]
  },
]);

