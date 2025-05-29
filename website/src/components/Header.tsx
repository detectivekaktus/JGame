import { Link } from "react-router-dom";
import { useContext, useEffect, useState } from "react"
import { useMediaQuery } from "../hooks/media_query";
import { MeContext } from "../context/MeProvider";
import "../css/Header.css"

export function Header() {
  const { me } = useContext(MeContext);
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
          <Link to={me ? "/main" : "/"}><img src="/JGame.svg" alt="JGame logo" /></Link>
          <button className="primary-nav-mobile-toggle" aria-controls="primary-navigation" aria-expanded={isMenuOpen && isSmall ? "true" : "false"} onClick={toggleMenu}>
            <img className="hamburger" src="/hamburger.png" alt="" aria-hidden="true" />
            <span className="sr-only">Menu</span>
          </button>
          {
            !me ?
          <nav className={`primary-nav ${isMenuOpen && isSmall ? "opened" : ""}`}>
            <ul aria-label="primary" className="primary-nav-list">
              <li><Link to="https://github.com/detectivekaktus/JGame">Source code</Link></li>
              <li><Link to="/packs">Quiz packs</Link></li>
              <li><Link to="/editor">Quiz editor</Link></li>
              <li><Link to="/auth/signup">Sign up</Link></li>
            </ul>
          </nav>
            :
          <nav className={`primary-nav ${isMenuOpen && isSmall ? "opened" : ""}`}>
            <ul aria-label="primary" className="primary-nav-list">
              <li><Link to="/main">Rooms</Link></li>
              <li><Link to="/packs">Quiz packs</Link></li>
              <li><Link to="/editor">Quiz editor</Link></li>
              <li><Link to={`/profile/${me.id}`}>Profile</Link></li>
            </ul>
          </nav>
          }
        </div>
      </div>
    </header>
  )
}
