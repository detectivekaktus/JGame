import { Link } from "react-router-dom";
import { Footer } from "../components/Footer"
import "../css/Form.css"

export function LoginPage() {
  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <div className="auth-two-col">
          <div className="auth-two-col-title">
            <h1>Welcome back! Hope you're about to have a good playing session!</h1>
          </div>
          <div className="auth-two-col-form">
            <h2>Log in</h2>
            <form>
              <div className="form-entry">
                <label htmlFor="email">Email</label>
                <input required id="email" name="email" type="email" />
              </div>
              <div className="form-entry">
                <label htmlFor="password">Password</label>
                <input required id="password" name="password" type="password" />
              </div>
              <button type="submit" className="button stretch">Log in</button>
            </form>
            <p>Don't have an account?  <Link className="fg-accent-600 underlined" to={"/auth/signup"}>Create one</Link></p>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
}

