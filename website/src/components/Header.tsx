import { Link } from "react-router-dom";
import { useEffect, useState } from "react"
import { useMediaQuery } from "../hooks/media_query";
import { Me } from "../types/user";
import { BASE_API_URL } from "../utils/consts";
import "../css/Header.css"

export function Header() {
  const [user, setUser] = useState<Me>({} as Me);
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const isSmall = useMediaQuery("(max-width: 50em)");

  const toggleMenu = () => {
    document.body.style.overflowY = isMenuOpen ? "visible" : "hidden";
    setIsMenuOpen(!isMenuOpen);
 }

  useEffect(() => {
    fetch(`${BASE_API_URL}/users/me`, { credentials: "include" })
      .then(res => {
        if (!res.ok) {
          console.error(`There was an error: got response ${res.status}.`);
          setUser({} as Me);
          return;
        }
        return res.json()
      })
      .then(data => setUser(data as Me))
      .catch(err => console.error(err))
  }, [])

  useEffect(() => {
    if (!isSmall && isMenuOpen)
      document.body.style.overflowY = "visible";
  }, [isSmall]);

  return (
    <header className="primary-header">
      <div className="container">
        <div className="nav-wrapper">
          <Link to={Object.keys(user).length === 0 ? "/" : "/main"}><img src="/JGame.svg" alt="JGame logo" /></Link>
          <button className="primary-nav-mobile-toggle" aria-controls="primary-navigation" aria-expanded={isMenuOpen && isSmall ? "true" : "false"} onClick={toggleMenu}>
            <img className="hamburger" src="/hamburger.png" alt="" aria-hidden="true" />
            <span className="sr-only">Menu</span>
          </button>
          {
            Object.keys(user).length === 0 ?
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
              <li><Link to={`/profile/${user.id}`}>Profile</Link></li>
            </ul>
          </nav>
          }
        </div>
      </div>
    </header>
  )
}
