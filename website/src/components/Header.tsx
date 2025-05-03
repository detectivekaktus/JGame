import { useState } from "react"
import "../css/Header.css"

// TODO: Update state of overlowY on body when the user goes beyond
// 40em with the navigation menu opened.
function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => {
    document.body.style.overflowY = isMenuOpen ? "visible" : "hidden";
    setIsMenuOpen(!isMenuOpen);
  }

  return (
    <header className="primary-header">
      <div className="container">
        <div className="nav-wrapper">
          <a href="#"><img src="/JGame.svg" alt="JGame logo" /></a>
          <button className="primary-nav-mobile-toggle" aria-controls="primary-navigation" aria-expanded={isMenuOpen ? "true" : "false"} onClick={toggleMenu}>
            <img className="hamburger" src="hamburger.png" alt="" aria-hidden="true" />
            <span className="sr-only">Menu</span>
          </button>
          <nav className={`primary-nav ${isMenuOpen ? "opened" : ""}`}>
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
