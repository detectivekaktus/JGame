import { useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import "../css/NotFoundPage.css"

export function NotFoundPage() {
  const navigate = useNavigate();
  return (
    <>
      <Header />
      <div className="content">
        <h1>404 Not Found</h1>
        <p>It looks like you're lost 🤔</p>
        <button className="button" onClick={() => navigate(-1)}>Go back</button>
      </div>
      <Footer />
    </>
  );
}
