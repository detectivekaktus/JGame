import "../css/RoomCard.css"

type RoomCardProps = {
  name: string;
  curUsers: number;
  maxUsers: number;
}

export function RoomCard({ name, curUsers, maxUsers }: RoomCardProps) {
  return (
    <button className="button room stretch">
      <h3>{name}</h3>
      <p>{curUsers}/{maxUsers}</p>
    </button>
  )
}
