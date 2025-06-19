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
import { WSActionType, WSMessage, WSQuestion, WSUser, WSUserRole } from "../utils/websocket";

type RoomParams = {
  id: string,
}

export function RoomPage() {
  const { id } = useParams<RoomParams>();

  const ws = useRef<WebSocket | null>(null);

  const [started, setStarted] = useState(false);
  const [role, setRole] = useState<WSUserRole>(WSUserRole.PLAYER);
  const [users, setUsers] = useState<WSUser[]>([])

  const [question, setQuestion] = useState<WSQuestion | null>(null);

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

    if (
      !ws.current ||
      ws.current.readyState === WebSocket.CLOSED ||
      ws.current.readyState === WebSocket.CLOSING
    )
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
      const msg: WSMessage = JSON.parse(e.data);
      console.debug(e.data)

      switch (msg.type) {
        case WSActionType.JOINED_ROOM: {
          setRole(msg.payload["role"]);
        } break;

        case WSActionType.USERS_LIST: {
          setUsers(msg.payload["users"]);
        } break;

        case WSActionType.GAME_STARTED: {
          setStarted(true);
        } break;

        case WSActionType.GAME_STATE: {
          setStarted(msg.payload["started"]);
        } break;

        case WSActionType.QUESTION: {
          setQuestion(msg.payload["question"]);
        } break;

        case WSActionType.QUESTIONS_DONE: {
          ws.current?.send(JSON.stringify({
            type: WSActionType.LEAVE_ROOM,
            payload: {
              room_id: Number(id)
            }
          } as WSMessage));
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
    } as WSMessage));

    ws.current?.send(JSON.stringify({
      type: WSActionType.NEXT_QUESTION,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage));
  };

  const handleLeave = () => {
    ws.current?.send(JSON.stringify({
      type: WSActionType.LEAVE_ROOM,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage));
  };

  const handleNextQuestion = () => {
    ws.current?.send(JSON.stringify({
      type: WSActionType.NEXT_QUESTION,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage));
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
                  { users.map((user, key) => <UserCard key={key} name={user.name} id={user.id} score={null}/>) }
                </div>
              </div>
            :
              <div className="margin-top room-in-game">
                <div className="question-panel">
                  <div className="question">{question?.title}</div>
                  <ol className="question-answer-options">
                    { question?.answers.map((answer, key) => <li key={key}><Button stretch={true} dim={false}>{answer.text}</Button></li>) }
                  </ol>
                </div>
                <div className="leaderboard">
                  <h2>Leaderboard</h2>
                  { users.sort((a, b) => b.score - a.score).map((user, key) => <UserCard key={key} name={user.name} id={user.id} score={user.score}/>) }
                </div>
              </div>
          }
          <div className="room-options">
            <ol>
              <li><Button stretch={false} dim={false} onClick={handleLeave}>Leave</Button></li>
              { role === WSUserRole.OWNER && !started && <li><Button stretch={false} dim={false} onClick={handleStart}>Start</Button></li> }
              { role === WSUserRole.OWNER && started && <li><Button stretch={false} dim={false} onClick={handleNextQuestion}>Next question</Button></li> }
            </ol>
          </div>
        </div>
        <Footer />
      </div>
    </>
  );
}
