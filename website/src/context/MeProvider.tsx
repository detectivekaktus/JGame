import { createContext, useEffect, useState } from "react";
import { Me } from "../types/user";
import { BASE_API_URL } from "../utils/consts";

interface MeContextType {
  me: Me | null
  setMe: React.Dispatch<React.SetStateAction<Me | null>>
}

export const MeContext = createContext<MeContextType>({
  me: null,
  setMe: () => {}
});

export function MeProvider({ children }: { children: React.ReactNode }) {
  const [me, setMe] = useState<Me | null>(null);

  useEffect(() => {
    fetch(`${BASE_API_URL}/users/me`, { credentials: "include" })
      .then(res => {
        if (res.ok)
          return res.json();
        else if (res.status === 401 || res.status === 403)
          return null;
        else
          throw new Error(`Unexpected error during current user fetch: ${res.status}`);
      })
      .then(data => data ? setMe(data) : {})
      .catch(err => console.error(err));
  }, []);

  return (
    <MeContext.Provider value={{ me, setMe }}>
      {children}
    </MeContext.Provider>
  );
}
