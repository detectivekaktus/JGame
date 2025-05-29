import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { StatBadge, StatBadgeColor } from "../components/StatBadge";
import { BASE_API_URL } from "../utils/consts";
import { User } from "../types/user";
import { NotFoundPage } from "./NotFoundPage";
import "../css/ProfilePage.css"

type ProfileParams = {
  id: string;
}

export function ProfilePage() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [found, setFound] = useState(false);

  const { id } = useParams<ProfileParams>();

  useEffect(() => {
    fetch(`${BASE_API_URL}/users/${id}`)
      .then(res => {
        if (res.ok)
          return res.json();
        else if (res.status === 404) {
          setFound(false);
          return null;
        }
        throw new Error(`Unexpected error during user fetch: ${res.status}`);
      })
      .then(data => {
        if (data) {
          setUser(data);
          setFound(true);
        }
      })
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  }, []);

  if (loading)
    return <h1>Loading...</h1>

  if (!found)
    return <NotFoundPage />

  return (
    <div className="page-wrapper">
      <Header />
      <main className="margin-top margin-bottom container profile page">
        <div className="name-picture">
          <h1>{user?.name}</h1>
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
