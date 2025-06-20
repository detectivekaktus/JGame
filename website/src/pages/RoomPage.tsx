import { closeSocket, getSocket, onMessageType, sendMessage, WSActionType, WSMessage, WSQuestion, WSUser, WSUserRole } from "../utils/websocket";
import { useContext, useEffect, useState } from "react";
import { MeContext } from "../context/MeProvider";
import { useNavigate, useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { BASE_API_URL } from "../utils/consts";
import "../css/RoomPage.css"
import { NotFoundPage } from "./NotFoundPage";
import { UserCard } from "../components/UserCard";

type RoomParams = {
  id: string,
}

export function RoomPage() {
  const { id } = useParams<RoomParams>();

  const [started, setStarted] = useState(false);
  const [finished, setFinished] = useState(false);
  const [role, setRole] = useState<WSUserRole>(WSUserRole.PLAYER);
  const [users, setUsers] = useState<WSUser[]>([])

  const [question, setQuestion] = useState<WSQuestion | null>(null);
  const [answered, setAnswered] = useState(false);

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

    const socket = getSocket();
    socket.addEventListener("open", () => {
      socket.send(JSON.stringify({
        type: WSActionType.JOIN_ROOM,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage));

      socket.send(JSON.stringify({
        type: WSActionType.GET_USERS,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage));

      socket.send(JSON.stringify({
        type: WSActionType.GET_GAME_STATE,
        payload: {
          room_id: Number(id)
        }
      } as WSMessage));
    });

    onMessageType(WSActionType.JOINED_ROOM, (msg: WSMessage) => setRole(msg.payload["role"]));
    onMessageType(WSActionType.USERS_LIST, (msg: WSMessage) => setUsers(msg.payload["users"]));
    onMessageType(WSActionType.GAME_STARTED, () => setStarted(true));
    onMessageType(WSActionType.QUESTION, (msg: WSMessage) => { setQuestion(msg.payload["question"]); setAnswered(false); })
    onMessageType(WSActionType.QUESTIONS_DONE, () => setFinished(true));
    onMessageType(WSActionType.LEFT_ROOM, () => navigate(-1));
    onMessageType(WSActionType.ROOM_DELETED, () => navigate(-1));
    onMessageType(WSActionType.GAME_STATE, (msg: WSMessage) => {
      const msgStarted = msg.payload["started"];
      setStarted(msgStarted);
      setFinished(msg.payload["finished"]);
      if (msgStarted)
        setQuestion(msg.payload["question"]);
    });

    setLoading(false);
    return () => closeSocket();
  }, [me, loadingMe, id])

  const handleStart = () => {
    sendMessage({
      type: WSActionType.START_GAME,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage);

    sendMessage({
      type: WSActionType.NEXT_QUESTION,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage);
  };

  const handleLeave = () => {
    sendMessage({
      type: WSActionType.LEAVE_ROOM,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage);
  };

  const handleNextQuestion = () => {
    sendMessage({
      type: WSActionType.NEXT_QUESTION,
      payload: {
        room_id: Number(id)
      }
    } as WSMessage);
  };

  const submitAnswer = (index: number) => {
    console.debug("answered")
    setAnswered(true);
    sendMessage({
      type: WSActionType.ANSWER,
      payload: {
        room_id: Number(id),
        answer: index
      }
    } as WSMessage);
  };

  if (loadingMe || loading)
    return <LoadingPage />

  if (!found)
    return <NotFoundPage />

  if (finished)
    return (
      <>
        <div className="page-wrapper">
          <div className="container page">
            <h1 className="margin-top">Results</h1>
            <div className="leaderboard">
              <h2>Leaderboard</h2>
              { users.sort((a, b) => b.score - a.score).map((user, key) => <UserCard key={key} name={user.name} id={user.id} score={user.score}/>) }
            </div>
            <Button className="margin-top" stretch={false} dim={false} onClick={handleLeave}>Leave</Button>
          </div>
        </div>
        <Footer />
      </>
    );

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
                    { question?.answers.map((answer, key) =>
                      <li key={key}><Button stretch={true} dim={false} onClick={() => submitAnswer(key)} disabled={answered}>{answer.text}</Button></li>
                    ) }
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
