import { useParams } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { StatBadge, StatBadgeColor } from "../components/StatBadge";
import "../css/ProfilePage.css"

type ProfileParams = {
  id: string;
}

export function ProfilePage() {
  const { id } = useParams<ProfileParams>();

  return (
    <div className="page-wrapper">
      <Header />
      <main className="margin-top margin-bottom container profile page">
        <div className="name-picture">
          { /* Render profile picture */ }
          <h1>Profile name</h1>
        </div>
        <div className="profile-stats">
          <h2>Stats</h2>
          <ul>
            <li><StatBadge title="Matches played" progress="0" color={StatBadgeColor.DARK}/></li>
            <li><StatBadge title="Matches won" progress="0" /></li>
            <li> <StatBadge title="Quizzes created" progress="0" color={StatBadgeColor.LIGHT} /></li>
          </ul>
        </div>
      </main>
      <Footer />
    </div>
  );
}
