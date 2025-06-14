import { useContext, useEffect, useState } from "react";
import { MeContext } from "../context/MeProvider";
import { useNavigate, useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { BASE_API_URL, BASE_WS_URL } from "../utils/consts";
import "../css/RoomPage.css"
import { NotFoundPage } from "./NotFoundPage";
import { User } from "../types/user";
import { UserCard } from "../components/UserCard";

type RoomParams = {
  id: string,
}

export function RoomPage() {
  const { id } = useParams<RoomParams>();

  const [ws, setWs] = useState<WebSocket | null>(null);

  const [started, setStarted] = useState(false);
  const [role, setRole] = useState("player");
  const [users, setUsers] = useState<User[]>([])

  const { me, loadingMe } = useContext(MeContext);
  const [loading, setLoading] = useState(true);
  const [found, setFound] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me)
      navigate("/auth/login");

    const checkRoomAndConnect = async () => {
      const res = await fetch(`${BASE_API_URL}/rooms/${id}`)
      if (res.status === 404) {
        setFound(false);
        setLoading(false);
        return;
      }

      const socket = new WebSocket(BASE_WS_URL);
      setWs(socket);
      setLoading(false);
    };

    checkRoomAndConnect();

    return () => ws?.close()
  }, [me, loadingMe, id])

  if (loadingMe || loading)
    return <LoadingPage />

  if (!found)
    return <NotFoundPage />

  return (
    <>
      <div className="page-wrapper">
        <div className="container page">
          {
            !started ?
              <div className="margin-top room-stats">
                <h1>Players connected</h1>
                <div className="players">
                  { users.map((user, key) => <UserCard key={key} name={user.name} id={user.id}/>) }
                </div>
                {role === "owner" && <Button stretch={false} dim={false} onClick={() => setStarted(true)}>Start</Button>}
              </div>
            :
              <div className="margin-top question">

              </div>
          }
        </div>
        <Footer />
      </div>
    </>
  );
}
