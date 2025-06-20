import { BASE_WS_URL } from "./consts"

export enum WSActionType {
  JOIN_ROOM        = "join_room",
  JOINED_ROOM      = "joined_room",

  LEAVE_ROOM       = "leave_room",
  LEFT_ROOM        = "left_room",
  ROOM_DELETED     = "room_deleted",

  START_GAME       = "start_game",
  GAME_STARTED     = "game_started",

  GET_USERS        = "get_users",
  USERS_LIST       = "users_list",

  GET_GAME_STATE   = "get_game_state",
  GAME_STATE       = "game_state",

  NEXT_QUESTION    = "next_question",
  QUESTION         = "question",
  QUESTIONS_DONE   = "questions_done",

  ANSWER           = "answer",

  ERROR            = "error"
}

export interface WSMessage {
  type: WSActionType
  payload: any
}

export enum WSUserRole {
  OWNER = "owner",
  PLAYER = "player"
}

export interface WSUser {
  id:      number
  name:    string
  role:    string
  room_id: number
  score:   number
}

export interface WSAnswer {
  text:    string
  correct: boolean
}

export interface WSQuestion {
  title:     string
  image_url: string
  value:     number
  answers:   WSAnswer[]
}

let socket: WebSocket | null = null;
const handlers: Map<WSActionType, (msg: WSMessage) => void> = new Map();

export function getSocket(): WebSocket {
  if (!socket || socket.readyState === WebSocket.CLOSED) {
    socket = new WebSocket(BASE_WS_URL)

    socket.onmessage = (e) => {
      console.log(e.data);
      const msg: WSMessage = JSON.parse(e.data);
      const handler = handlers.get(msg.type);

      if (handler)
        handler(msg);
      else {
        console.warn(`No handler for '${msg.type}'`);
        console.warn(JSON.stringify(msg))
      }
    };

    socket.onclose = () => console.log("websocket died.");
    socket.onerror = (err) => console.error(err);
  }

  return socket;
}

export function sendMessage(msg: WSMessage) {
  const sock = getSocket();

  if (sock.readyState === WebSocket.OPEN)
    sock.send(JSON.stringify(msg));
  else if (sock.readyState === WebSocket.CONNECTING)
    sock.addEventListener("open", () => sock.send(JSON.stringify(msg)), { once: true });
}

export function onMessageType(type: WSActionType, handler: (msg: WSMessage) => void) {
  handlers.set(type, handler);
}

export function closeSocket() {
  if (socket) {
    socket.close();
    socket = null;
    handlers.clear();
  }
}
