import { useContext, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { StatBadge, StatBadgeColor } from "../components/StatBadge";
import { BASE_API_URL } from "../utils/consts";
import { User } from "../types/user";
import { MeContext } from "../context/MeProvider";
import { NotFoundPage } from "./NotFoundPage";
import { LoadingPage } from "./LoadingPage";
import "../css/ProfilePage.css"

type ProfileParams = {
  id: string;
}

export function ProfilePage() {
  const { me, loadingMe } = useContext(MeContext);

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [found, setFound] = useState(false);

  const { id } = useParams<ProfileParams>();
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me) {
      navigate("/auth/login");
      return;
    }

    if (me.id === Number(id)) {
      setFound(true);
      setLoading(false);
      setUser({
        id: me.id,
        name: me.name,
        matches_played: me.matches_played,
        matches_won: me.matches_won
      });
      return;
    }

    fetch(`${BASE_API_URL}/users/${id}`)
      .then(res => {
        if (res.ok)
          return res.json();
        else if (res.status === 404) {
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
  }, [id, me, loadingMe]);

  if (loadingMe || loading)
    return <LoadingPage />

  if (!found)
    return <NotFoundPage />

  return (
    <div className="page-wrapper">
      <Header />
      <main className="margin-top margin-bottom container profile page">
        <div className="name-picture">
          <h1>{ me?.id === Number(id) && <Link className="profile-settings" to="/profiles/settings"><img src="/settings.svg" alt="Settings icon" /></Link> }{user?.name}</h1>
        </div>
        <div className="profile-stats">
          <h2>Stats</h2>
          <ul>
            <li><StatBadge title="Matches played" progress={user?.matches_played.toString() || "0"} color={StatBadgeColor.DARK}/></li>
            <li><StatBadge title="Matches won" progress={user?.matches_won.toString() || "0"} /></li>
          </ul>
        </div>
      </main>
      <Footer />
    </div>
  );
}
