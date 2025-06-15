import { useContext, useEffect, useRef, useState } from "react";
import { MeContext } from "../context/MeProvider";
import { useNavigate, useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { BASE_API_URL, BASE_WS_URL } from "../utils/consts";
import "../css/RoomPage.css"
import { NotFoundPage } from "./NotFoundPage";
import { UserCard } from "../components/UserCard";
import { WSActionType, WSMessage, WSUser } from "../utils/websocket";
import { User } from "../types/user";

type RoomParams = {
  id: string,
}

export function RoomPage() {
  const { id } = useParams<RoomParams>();

  const ws = useRef<WebSocket | null>(null);

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

    const checkRoom = async () => {
      const res = await fetch(`${BASE_API_URL}/rooms/${id}`)
      if (res.status === 404) {
        setFound(false);
        return;
      }
    };
    checkRoom();

    if (!ws.current || ws.current.readyState === WebSocket.CLOSED || ws.current.readyState === WebSocket.CLOSING)
      ws.current = new WebSocket(BASE_WS_URL);

    ws.current.onopen = () => {
      ws.current?.send(JSON.stringify({
        type: WSActionType.JOIN_ROOM,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage));

      ws.current?.send(JSON.stringify({
        type: WSActionType.GET_USERS,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage))

      ws.current?.send(JSON.stringify({
        type: WSActionType.GET_GAME_STATE,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage))
    };
    ws.current.onmessage = (e) => {
      const msg: WSMessage = JSON.parse(e.data)
      console.debug(e.data)

      switch (msg.type) {
        case WSActionType.JOINED_ROOM: {
          setRole(msg.payload["role"])
        } break;

        case WSActionType.USERS_LIST: {
          const wsUsers: WSUser[] = msg.payload["users"]
          Promise.all(
            wsUsers.map((user) =>
              fetch(`${BASE_API_URL}/users/${user.id}`).then((res) => res.json())
            )
          ).then((fetchedUsers) => {
              setUsers(fetchedUsers);
            });
        } break;

        case WSActionType.GAME_STARTED: {
          setStarted(true);
        } break;

        case WSActionType.GAME_STATE: {
          setStarted(msg.payload["started"]);
        } break;

        case WSActionType.LEFT_ROOM:
        case WSActionType.ROOM_DELETED: {
          navigate(-1);
        } break;

        default: {
          console.error(`encountered unknown action type ${msg.type}: ${JSON.stringify(msg.payload)}`)
          ws.current?.close();
        } break;
      }
    };

    setLoading(false);
    const wsCurrent = ws.current;
    return () => wsCurrent.close()
  }, [me, loadingMe, id])

  const handleStart = () => {
    ws.current?.send(JSON.stringify({
      type: WSActionType.START_GAME,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage))
  };

  const handleLeave = () => {
    ws.current?.send(JSON.stringify({
      type: WSActionType.LEAVE_ROOM,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage))
  };

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
                <div className="room-options">
                  <ol>
                    <li><Button stretch={false} dim={false} onClick={handleLeave}>Leave</Button></li>
                    {role === "owner" && <li><Button stretch={false} dim={false} onClick={handleStart}>Start</Button></li>}
                  </ol>
                </div>
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
