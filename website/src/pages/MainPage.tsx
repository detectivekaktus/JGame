import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import "../css/MainPage.css"

export function MainPage() {
  const { me, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me)
      navigate("/");
  }, [me, loadingMe])

  if (loadingMe)
    return <LoadingPage />

  return (
    <div className="page-wrapper">
      <Header />
      <main className="page">
        <div className="margin-top margin-bottom container main-menu">
          <div className="menu-rooms">
            <h2>Rooms</h2>
            {
              true ? <h2>There are no rooms</h2> : 
                <ul className="menu-rooms-list">
                  { /* Fetch and render rooms */ }
                </ul>
            }
          </div>
          <div className="menu-nav">
            <div className="bg-accent-600 menu-nav-room-desc">
              {
                true ? <h2>Select a room to see its details</h2> :
                  <div className="menu-nav-room">
                    <h2>Room name</h2>
                    <ul>
                      <li>Room type: </li>
                      <li>Players in room: </li>
                      <li>Quiz pack: </li>
                      <li><Button stretch={true} dim={true} >Join</Button></li>
                    </ul>
                  </div>
              }
            </div>
            <Button stretch={true} dim={false} >Create room</Button>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
