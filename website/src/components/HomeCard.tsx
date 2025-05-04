import "../css/HomeCard.css"

type HomeCardProps = {
  title: string;
  description: string;
  invert?: boolean;
  children: React.ReactNode;
}

export function HomeCard({ title, description, invert, children }: HomeCardProps) {
  return (
    <div className="card" datatype={ invert ? "inverted" : ""}>
      <div className="card-description">
        <h2>{title}</h2>
        <p>{description}</p>
      </div>
      <div className="card-content">
        {children}
      </div>
    </div>
  )
}
