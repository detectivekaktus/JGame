import { RouterProvider} from "react-router-dom"
import { APP_ROUTER } from "./router";

export function App() {
  const router = APP_ROUTER;
  return <RouterProvider router={router} />;
}

