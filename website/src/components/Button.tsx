type ButtonProps = {
  stretch: boolean,
  dim: boolean,
  type?: "button" | "submit" | "reset",
  disabled?: boolean,

  onClick?: () => void,
  children: any
}

export function Button({ stretch, dim, type, disabled, onClick, children }: ButtonProps) {
  return (
    <button
      onClick={onClick}
      className={`button ${stretch ? "stretch" : ""}`}
      datatype={dim ? "dim" : ""}
      type={type}
      disabled={disabled}>
      {children}
    </button>
  );
}
