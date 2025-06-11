import "../css/RoomCard.css"

export type Room = {
  room_id: number,
  name: string,
  pack_id: number,
  user_id: number,
  users: number[],
  current_users: number,
  max_users: number,
}

type RoomCardProps = {
  name: string,
  curUsers: number,
  maxUsers: number,
  onClick: () => void
}

export function RoomCard({ name, curUsers, maxUsers, onClick }: RoomCardProps) {
  return (
    <button onClick={onClick} className="button room stretch">
      <h3>{name}</h3>
      <p>{curUsers}/{maxUsers}</p>
    </button>
  )
}
