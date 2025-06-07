import { useEffect, useState } from "react";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { Pack, PackCard } from "../components/PackCard";
import { BASE_API_URL } from "../utils/consts";
import "../css/PacksPage.css"
import { LoadingPage } from "./LoadingPage";

export function PacksPage() {
  const [packs, setPacks] = useState<Pack[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${BASE_API_URL}/packs`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setPacks(data))
      .catch(err => console.error(err))
    .finally(() => setLoading(false))
  }, []);

  if (loading)
    return <LoadingPage />

  return (
    <>
      <div className="page-wrapper">
        <Header />
        <div className="page">
          <div className="margin-top packs-menu container">
            <h2>Packs</h2>
            <div className="packs">
              <ul className="packs-list">
                { 
                  packs.length === 0 ?
                    <h2>No packs currently exist</h2>
                  :
                    packs.map((pack) => <PackCard key={pack.id} pack={pack} /> )
                }
              </ul>
            </div>
          </div>
        </div>
        <Footer />
      </div>
    </>
  );
}
