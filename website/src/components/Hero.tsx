import "../css/Hero.css"

function Hero() {
  return (
    <section className="margin-top container hero">
      <div className="hero-left-wrapper">
        <div className="hero-description">
          <h1>Custom quizzes right in your browser</h1>
          <p>Enjoy puzzles with your friends anywhere at any time</p>
        </div>
        <div className="hero-buttons">
          <button className="button">Sign up</button>
          <button className="button" datatype="dim">Log in</button>
        </div>
      </div>
      <div className="hero-example">
        <h2>When was the declaration of independence signed?</h2>
        <ul className="hero-example-options">
          <li><button className="button" datatype="stretch">August 2, 1776</button></li>
          <li><button className="button" datatype="stretch">September 1, 1781</button></li>
          <li><button className="button" datatype="stretch">July 2, 1770</button></li>
          <li><button className="button" datatype="stretch">April 19, 1775</button></li>
        </ul>
      </div>
    </section>
  )
}

export default Hero
