import "../css/StatBadge.css"

export enum StatBadgeColor {
  LIGHT,
  NORMAL,
  DARK
}

type StatBadgeProps = {
  title: string;
  progress: string;
  color?: StatBadgeColor
}

function resolveColor(color: StatBadgeColor | undefined): string {
  switch (color) {
    case StatBadgeColor.LIGHT:
      return "bg-accent-300";
    case StatBadgeColor.NORMAL:
      return "bg-accent-600";
    case StatBadgeColor.DARK:
      return "bg-accent-700";
    default:
      return "bg-accent-600";
  }
}

export function StatBadge({ title, progress, color }: StatBadgeProps) {
  return (
    <div className={`stat-badge ${resolveColor(color)}`}>
      <h3>{title}</h3>
      <p>{progress}</p>
    </div>
  )
}
