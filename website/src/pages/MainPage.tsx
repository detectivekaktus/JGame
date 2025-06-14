import React, { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { Room, RoomCard, RoomRequestForm } from "../components/RoomCard";
import { BASE_API_URL } from "../utils/consts";
import { Spinner } from "../components/Spinner";
import { Search } from "../components/Search";
import "../css/MainPage.css"

export function MainPage() {
  const [query, setQuery] = useState("");
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);

  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
  const [selectedRoomOwnerName, setSelectedRoomOwnerName] = useState("unknown");
  const [selectedRoomPackName, setSelectedRoomPackName] = useState("unknown");

  const [createRoomMenuDisplayed, setCreateRoomMenuDisplayed] = useState(false);
  const [createRoomErrors, setCreateRoomErorrs] = useState<Record<string, string>>({});

  const { me, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me)
      navigate("/");

    fetch(`${BASE_API_URL}/rooms`)
      .then(res => res.ok ? res.json() : [])
      .then(data => setRooms(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  }, [me, loadingMe])

  useEffect(() => {
    if (!selectedRoom)
      return;

    fetch(`${BASE_API_URL}/users/${selectedRoom?.user_id}`)
      .then(res => res.json())
      .then(data => setSelectedRoomOwnerName(data["name"]))
      .catch(err => console.error(err));

    fetch(`${BASE_API_URL}/packs/${selectedRoom?.pack_id}`)
      .then(res => res.json())
      .then(data => setSelectedRoomPackName(data["name"]))
      .catch(err => console.error(err));
  }, [selectedRoom])

  const handleQuery = () => {
    setLoading(true);
    fetch(`${BASE_API_URL}/rooms?name=${query}`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setRooms(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  };

  const handleCreateRoom = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const target = e.currentTarget;
    const data = new FormData(target);
    const errors: Record<string, string> = {};
    
    const name = data.get("name");
    const pack_id = data.get("pack_id");
    const password = data.get("password");
    if (typeof name !== "string" || typeof pack_id !== "string" || typeof password != "string") {
      errors.format = "Invalid format.";
      setCreateRoomErorrs(errors);
      return;
    }

    if (name.length < 4 || name.length > 32) {
      errors.name = "Name must be 4 to 32 characters long.";
      setCreateRoomErorrs(errors);
      return;
    }

    if (password.length > 32) {
      errors.password = "Password max length is 32 characters.";
      setCreateRoomErorrs(errors);
      return;
    }

    let res = await fetch(`${BASE_API_URL}/packs/${pack_id}`)
    if (!res.ok) {
      errors.pack_id = "This pack doesn't exist.";
      setCreateRoomErorrs(errors);
      return;
    }

    const room: RoomRequestForm = { name, pack_id: Number(pack_id), password: password === "" ? "PASSWORD_UNSET" : password}
    res = await fetch(`${BASE_API_URL}/rooms`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify(room)
    })

    const body = await res.json();
    if (!res.ok) {
      errors.req = body["message"];
      setCreateRoomErorrs(errors);
      return;
    }

    navigate(`/room/${body["room_id"]}`);
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
                    rooms.map((room, key) =>
                      <RoomCard key={key} onClick={() => setSelectedRoom(room)} name={room.name} curUsers={room.current_users} maxUsers={room.max_users}/>)
              }
            </ul>
          </div>
          <div className="menu-nav">
            <div className="bg-accent-600 menu-nav-room-desc">
              {
                !selectedRoom ?
                  <h2>Select a room to see its details</h2>
                :
                  <div className="menu-nav-room">
                    <h2>{selectedRoom.name}</h2>
                    <ul>
                      <li>Created by: <strong>{selectedRoomOwnerName}</strong></li>
                      <li>Quiz pack: <strong>{selectedRoomPackName}</strong></li>
                      <li>Players in room: <strong>{selectedRoom.current_users}</strong></li>
                      <li>Maximum players in room: <strong>{selectedRoom.max_users}</strong></li>
                      <li><Button stretch={true} dim={true} >Join</Button></li>
                    </ul>
                  </div>
              }
            </div>
            <Button stretch={true} dim={false} onClick={() => setCreateRoomMenuDisplayed(!createRoomMenuDisplayed)}>Create room</Button>
            {
              createRoomMenuDisplayed && (
                <div className="create-room-popup">
                  <h2>Create room</h2>
                  { createRoomErrors.req && <div className="form-error">{createRoomErrors.req}</div> }
                  { createRoomErrors.format && <div className="form-error">{createRoomErrors.format}</div> }
                  <form onSubmit={handleCreateRoom} noValidate>
                    <div className="form-entry">
                      <label htmlFor="name">Room name</label>
                      <input required id="name" name="name" type="text" defaultValue={"Room"} />
                      { createRoomErrors.name && <div className="form-entry-error">{createRoomErrors.name}</div> }
                    </div>
                    <div className="form-entry">
                      <label htmlFor="name">Pack id</label>
                      <input required id="pack_id" name="pack_id" type="number" min={1} defaultValue={1} />
                      { createRoomErrors.pack_id && <div className="form-entry-error">{createRoomErrors.pack_id}</div> }
                    </div>
                    <div className="form-entry">
                      <label htmlFor="name">Password</label>
                      <input id="password" name="password" type="text" />
                      { createRoomErrors.password && <div className="form-entry-error">{createRoomErrors.password}</div> }
                    </div>
                    <Button stretch={true} dim={false} type="submit">Create</Button>
                  </form>
                </div>
              )
            }
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
