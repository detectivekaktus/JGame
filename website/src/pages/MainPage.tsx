import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { Room, RoomCard } from "../components/RoomCard";
import { BASE_API_URL } from "../utils/consts";
import { Spinner } from "../components/Spinner";
import { Search } from "../components/Search";
import "../css/MainPage.css"

export function MainPage() {
  const [query, setQuery] = useState("");
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);

  const { me, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me)
      navigate("/");

    fetch(`${BASE_API_URL}/rooms`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setRooms(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  }, [me, loadingMe])

  const handleQuery = () => {
    setLoading(true);
    fetch(`${BASE_API_URL}/rooms?name=${query}`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setRooms(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  };

  if (loadingMe)
    return <LoadingPage />

  return (
    <div className="page-wrapper">
      <Header />
      <main className="page">
        <Search placeholder="Type room name..." setQuery={setQuery} handleQuery={handleQuery} />
        <div className="margin-top margin-bottom container main-menu">
          <div className="menu-rooms">
            <h2>Rooms</h2>
            <ul className="menu-rooms-list">
              {
                loading ?
                  <Spinner />
                : 
                  !rooms || rooms?.length === 0 ?
                    <h2>There are no rooms</h2>
                  :
                    rooms.map((room) => <RoomCard name={room.name} curUsers={room.current_users} maxUsers={room.max_users}/>)
              }
            </ul>
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
