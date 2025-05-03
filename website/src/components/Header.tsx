import { useEffect, useState } from "react"
import "../css/Header.css"

function useMediaQuery(query: string): boolean {
  const [matches, setMatches] = useState(() => window.matchMedia(query).matches);

  useEffect(() => {
    const mql: MediaQueryList = window.matchMedia(query);
    const handler = (event: MediaQueryListEvent) => setMatches(event.matches);
    mql.addEventListener("change", handler);
    return () => mql.removeEventListener("change", handler);
  }, [query]);

  return matches;
}

function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const isSmall = useMediaQuery("(max-width: 50em)");

  const toggleMenu = () => {
    document.body.style.overflowY = isMenuOpen ? "visible" : "hidden";
    setIsMenuOpen(!isMenuOpen);
  }

  useEffect(() => {
    if (!isSmall && isMenuOpen)
      document.body.style.overflowY = "visible";
  }, [isSmall]);

  return (
    <header className="primary-header">
      <div className="container">
        <div className="nav-wrapper">
          <a href="#"><img src="/JGame.svg" alt="JGame logo" /></a>
          <button className="primary-nav-mobile-toggle" aria-controls="primary-navigation" aria-expanded={isMenuOpen && isSmall ? "true" : "false"} onClick={toggleMenu}>
            <img className="hamburger" src="hamburger.png" alt="" aria-hidden="true" />
            <span className="sr-only">Menu</span>
          </button>
          <nav className={`primary-nav ${isMenuOpen && isSmall ? "opened" : ""}`}>
            <ul aria-label="primary" className="primary-nav-list">
              <li><a href="https://github.com/detectivekaktus/JGame">Source code</a></li>
              <li><a href="/packs">Quiz packs</a></li>
              <li><a href="/editor">Quiz editor</a></li>
              <li><a href="/signup">Sign up</a></li>
            </ul>
          </nav>
        </div>
      </div>
    </header>
  )
}

export default Header
