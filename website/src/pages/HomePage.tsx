import { useContext, useEffect } from "react"
import { Link, useNavigate } from "react-router-dom"
import { Footer } from "../components/Footer"
import { Header } from "../components/Header"
import { HomeCard } from "../components/HomeCard"
import { RoomCard } from "../components/RoomCard"
import { StatBadge, StatBadgeColor } from "../components/StatBadge"
import { MeContext } from "../context/MeProvider"
import { LoadingPage } from "./LoadingPage"
import { Button } from "../components/Button"
import "../css/Home.css"

export function HomePage() {
  const { me, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (me)
      navigate("/main");
  }, [me, loadingMe]);

  if (loadingMe)
    return <LoadingPage />

  return (
    <>
      <Header />
      <main>
        <section className="margin-top container hero">
          <div className="hero-left-wrapper">
            <div className="hero-description">
              <h1>Custom quizzes right in your browser</h1>
              <p>Enjoy puzzles with your friends anywhere at any time</p>
            </div>
            <div className="hero-buttons">
              <Link to="/auth/signup"><Button stretch={false} dim={false}>Sign up</Button></Link>
              <Link to="/auth/login"><Button stretch={false} dim={true}>Log in</Button></Link>
            </div>
          </div>
          <div className="hero-example">
            <h2>When was the declaration of independence signed?</h2>
            <ul className="hero-example-options">
              <li><Button stretch={true} dim={false} >August 2, 1776</Button></li>
              <li><Button stretch={true} dim={false} >September 1, 1781</Button></li>
              <li><Button stretch={true} dim={false} >July 2, 1770</Button></li>
              <li><Button stretch={true} dim={false} >April 19, 1775</Button></li>
            </ul>
          </div>
        </section>
        <section className="margin-top container cards">
          <HomeCard title="Join other people or create your own room"
            description="Bring your friends or strangers into your own room to play or join an already existing room to play the game.">
            <div className="rooms">
              <RoomCard name="John's Hills" curUsers={2} maxUsers={16}/>
              <RoomCard name="Fun together" curUsers={9} maxUsers={16}/>
              <RoomCard name="Intellectual battlefield" curUsers={5} maxUsers={16}/>
            </div>
          </HomeCard>
          <HomeCard title="Keep track of your progress on your profile"
            description="In your profile you can see your stats such as total matches played, win percentage, achievements and much more!"
            invert={true}>
            {/* TODO: I still don't know how to make stat badges be a perfect square box that stick to the left*/}
            <div className="stats">
              <StatBadge title="Total matches" progress="49" color={StatBadgeColor.LIGHT}/>
              <StatBadge title="Total won" progress="27"/>
              <StatBadge title="Quizzes created" progress="5" color={StatBadgeColor.DARK}/>
            </div>
          </HomeCard>
          <HomeCard title="Create your custom quiz packs"
            description="Done with other people’s quizzes or want to create a specific pack? You can always do it with the editor">
            <div className="card-question-wrapper">
              <img className="card-question" src="/home-card-question.png" alt="Question: The declaration of independence was signed on August 2, 1776" />
            </div>
          </HomeCard>
        </section>
        <section className="margin-top bg-accent-600 cta">
          <div className="container">
            <p>Sign up today to get the unforgettable quiz experience</p>
            <Link to="/auth/signup"><Button stretch={true} dim={true}>Sign up</Button></Link>
          </div>
        </section>
      </main>
      <Footer />
    </>
  )
}

