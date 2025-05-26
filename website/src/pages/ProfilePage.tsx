import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { StatBadge, StatBadgeColor } from "../components/StatBadge";
import { BASE_API_URL } from "../utils/consts";
import { User } from "../types/user";
import "../css/ProfilePage.css"
import { NotFoundPage } from "./NotFoundPage";

type ProfileParams = {
  id: string;
}

export function ProfilePage() {
  const [user, setUser] = useState<User>({} as User);
  const [found, setFound] = useState(true);
  const { id } = useParams<ProfileParams>();
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${BASE_API_URL}/users/me`, { credentials: "include" })
      .then(res => {
        if (res.status === 401 || res.status === 403)
          navigate("/")
      })
      .catch(err => console.error(err))

    fetch(`${BASE_API_URL}/users/${id}`)
      .then(res => {
        if (res.status == 404)
          setFound(false)
        return res.json()
      })
      .then(data => setUser(data as User))
  }, []);

  if (!found)
    return <NotFoundPage />

  return (
    <div className="page-wrapper">
      <Header />
      <main className="margin-top margin-bottom container profile page">
        <div className="name-picture">
          <h1>{user.name}</h1>
        </div>
        <div className="profile-stats">
          <h2>Stats</h2>
          <ul>
            <li><StatBadge title="Matches played" progress="0" color={StatBadgeColor.DARK}/></li>
            <li><StatBadge title="Matches won" progress="0" /></li>
            <li><StatBadge title="Quizzes created" progress="0" color={StatBadgeColor.LIGHT} /></li>
          </ul>
        </div>
      </main>
      <Footer />
    </div>
  );
}
