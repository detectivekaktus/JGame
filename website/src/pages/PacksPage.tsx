import { useEffect, useState } from "react";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import { Pack, PackCard } from "../components/PackCard";
import { BASE_API_URL } from "../utils/consts";
import { Spinner } from "../components/Spinner";
import { Search } from "../components/Search";
import "../css/PacksPage.css"

export function PacksPage() {
  const [query, setQuery] = useState("");
  const [packs, setPacks] = useState<Pack[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${BASE_API_URL}/packs`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setPacks(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  }, []);

  const handleQuery = () => {
    setLoading(true);
    fetch(`${BASE_API_URL}/packs?name=${query}`)
      .then(res => res.status === 200 ? res.json() : [])
      .then(data => setPacks(data))
      .catch(err => console.error(err))
      .finally(() => setLoading(false));
  };

  return (
    <>
      <div className="page-wrapper">
        <Header />
        <div className="page">
          <Search placeholder="Type pack name..." setQuery={setQuery} handleQuery={handleQuery} />
          <div className="margin-top packs-menu container">
            <h2>Packs</h2>
            {
              loading ?
                <Spinner />
              :
                <div className="packs">
                  <ul className="packs-list">
                    { 
                      !packs || packs?.length === 0 ?
                        <h2>No packs found</h2>
                        :
                        packs.map((pack) => <PackCard key={pack.id} pack={pack} /> )
                    }
                  </ul>
                </div>
            }
          </div>
        </div>
        <Footer />
      </div>
    </>
  );
}
