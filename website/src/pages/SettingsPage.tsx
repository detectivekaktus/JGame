import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import "../css/SettingsPage.css"
import { Button } from "../components/Button";
import { Footer } from "../components/Footer";

export function SettingsPage() {
  const { me, setMe, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me) {
      navigate("/auth/login");
      return;
    }
  }, [me, loadingMe]);

  if (loadingMe)
    return <LoadingPage />

  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <h1>Settings</h1>
        <hr />
        <dl>
          <div className="dl-section">
            <dt>Public information</dt>
            <dd>
              <p>Your name</p>
              <input type="text" value={me?.name}/>
            </dd>
          </div>
          <div className="dl-section">
            <dt>Private information</dt>
            <dd>
              <p>Your email</p>
              <input type="text" value={me?.email} />
            </dd>
          </div>
          <div className="dl-section">
            <dt>Password</dt>
            <dd>
              <p>Current password</p>
              <input type="text" />
              <p>New password</p>
              <input type="text" />
              <p>Re-type new password</p>
              <input type="text" />
            </dd>
          </div>
        </dl>
        <div className="button-options">
          <Button stretch={true} dim={false} onClick={() => navigate(-1)}>← Go back</Button>
          <Button className="save-changes-button" stretch={true} dim={false}>Save changes</Button>
        </div>
        <hr />
        <Button stretch={true} dim={false}>Delete account</Button>
      </div>
      <Footer />
    </div>
  );
}
