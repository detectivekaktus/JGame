export enum WSActionType {
  JOIN_ROOM    = "join_room",
  JOINED_ROOM  = "joined_room",

  LEAVE_ROOM   = "leave_room",
  LEFT_ROOM    = "left_room",

  START_GAME   = "start_game",
  GAME_STARTED = "game_started",

  GET_USERS    = "get_users",
  USERS_LIST   = "users_list",

  ERROR        = "error"
}

export interface WSMessage {
  type: WSActionType
  payload: any
}

