import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { LoginPage } from "./pages/LoginPage";
import { SignupPage } from "./pages/SignupPage";
import { MainPage } from "./pages/MainPage";
import { ProfilePage } from "./pages/ProfilePage";
import { SettingsPage } from "./pages/SettingsPage";
import { PacksPage } from "./pages/PacksPage";
import { RoomPage } from "./pages/RoomPage";

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
        path: "profiles",
        children: [
          { index: true, Component: NotFoundPage },
          { path: ":id", Component: ProfilePage },
          { path: "settings", Component: SettingsPage }
        ]
      },
      { path: "packs", Component: PacksPage },
      {
        path: "rooms",
        children: [
          { index: true, Component: NotFoundPage },
          { path: ":id", Component: RoomPage }
        ]
      },
      { path: "*", Component: NotFoundPage }
    ]
  },
]);

