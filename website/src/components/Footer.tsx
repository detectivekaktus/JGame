import "../css/Footer.css"

export function Footer() {
  return (
    <footer className="bg-neutral-700 main-footer">
      <div className="container">
        <div className="copyright">
          <img src="/JGame.svg" alt="JGame logo" />
          <p>Made with ❤️ by Artiom Astashonak<br />
            Copyright © 2025. All rights reserved.</p>
        </div>
        <nav className="main-footer-nav">
          <ul className="main-footer-list">
            <a href="https://github.com/detectivekaktus/JGame"><li>Source code</li></a>
            <a href="/packs"><li>Quiz packs</li></a>
            <a href="/editor"><li>Quiz editor</li></a>
            <a href="/signup"><li>Sign up</li></a>
            <a href="/login"><li>Log in</li></a>
          </ul>
        </nav>
      </div>
    </footer>
  )
}

