type ButtonProps = {
  stretch: boolean,
  dim: boolean,
  className?: string,
  type?: "button" | "submit" | "reset",
  disabled?: boolean,

  onClick?: () => void,
  children: any
}

export function Button({ stretch, dim, type, className, disabled, onClick, children }: ButtonProps) {
  return (
    <button
      onClick={onClick}
      className={`button ${stretch ? "stretch" : ""} ${className ? className : ""}`}
      datatype={dim ? "dim" : ""}
      type={type ? type : "button"}
      disabled={disabled}>
      {children}
    </button>
  );
}
