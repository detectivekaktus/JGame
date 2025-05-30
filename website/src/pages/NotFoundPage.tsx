import { useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer";
import { Header } from "../components/Header";
import "../css/NotFoundPage.css"
import { Button } from "../components/Button";

export function NotFoundPage() {
  const navigate = useNavigate();
  return (
    <div className="page-wrapper">
      <Header />
      <div className="page content">
        <h1>404 Not Found</h1>
        <p>It looks like you're lost 🤔</p>
        <Button stretch={false} dim={false} onClick={() => navigate(-1)}>Go back</Button>
      </div>
      <Footer />
    </div>
  );
}
