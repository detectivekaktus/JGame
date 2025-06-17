import { Link } from "react-router-dom";
import "../css/UserCard.css"

type UserCardProps = {
  id: number,
  name: string,
  score: number | null
}

export function UserCard({ id, name, score }: UserCardProps) {
  return (
    <div className="user-card">
      <Link to={`/profiles/${id}`} target="_blank" rel="noopener noreferrer">{name}</Link>
      {score !== null && score}
    </div>
  );
}
