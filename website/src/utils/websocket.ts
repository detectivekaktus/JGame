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
