import { Link } from "react-router-dom"
import "../css/Footer.css"

export function Footer() {
  return (
    <footer className="bg-neutral-700 main-footer">
      <div className="container">
        <div className="copyright">
          <Link to="/"><img src="/JGame.svg" alt="JGame logo" /></Link>
          <p>Made with ❤️ by Artiom Astashonak<br />
            Copyright © 2025. All rights reserved.</p>
        </div>
        <nav className="main-footer-nav">
          <ul className="main-footer-list">
            <Link to="https://github.com/detectivekaktus/JGame"><li>Source code</li></Link>
            <Link to="/packs"><li>Quiz packs</li></Link>
            <Link to="/editor"><li>Quiz editor</li></Link>
            <Link to="/auth/signup"><li>Sign up</li></Link>
            <Link to="/auth/login"><li>Log in</li></Link>
          </ul>
        </nav>
      </div>
    </footer>
  )
}

