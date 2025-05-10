import { Link } from "react-router-dom";
import { Footer } from "../components/Footer"
import "../css/Form.css"

export function SignupPage() {
  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <div className="auth-two-col">
          <div className="auth-two-col-title">
            <h1>Good to see you! Let’s sign up, it’s quick and simple</h1>
          </div>
          <div className="auth-two-col-form">
            <h2>Sign up</h2>
            <form>
              <div className="form-entry">
                <label htmlFor="name">Display name</label>
                <input required id="name" name="name" type="text" />
              </div>
              <div className="form-entry">
                <label htmlFor="email">Email</label>
                <input required id="email" name="email" type="email" />
              </div>
              <div className="form-entry">
                <label htmlFor="password">Password</label>
                <input required id="password" name="password" type="password" />
              </div>
              <div className="form-entry">
                <label htmlFor="re-password">Repeat password</label>
                <input required id="re-password" name="re-password" type="password" />
              </div>
              <button type="submit" className="button stretch" >Sign up</button>
            </form>
            <p>Already have an account?  <Link className="fg-accent-600 underlined" to={"/auth/login"}>Log in</Link></p>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
}

