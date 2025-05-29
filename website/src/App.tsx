import { RouterProvider} from "react-router-dom"
import { MeProvider } from "./context/MeProvider";
import { APP_ROUTER } from "./router";

export function App() {
  const router = APP_ROUTER;
  return (
    <MeProvider>
      <RouterProvider router={router} />
    </MeProvider>
  );
}

