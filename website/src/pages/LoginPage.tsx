import { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer"
import { LoginForm } from "../types/user";
import { BASE_API_URL } from "../utils/consts";
import { MeContext } from "../context/MeProvider";
import "../css/Form.css"

export function LoginPage() {
  const { me, setMe } = useContext(MeContext);

  const [errors, setErrors] = useState<Record<string, string>>({});
  const navigate = useNavigate();

  useEffect(() => {
    if (me)
      navigate("/main");
  }, [me]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const target: HTMLFormElement = e.currentTarget;
    const data = new FormData(target);
    const userForm: LoginForm = {
      email: data.get("email") as string,
      password: data.get("password") as string,
    };
    const errors: Record<string, string> = {};

    // Copied from StackOverflow
    const emailRegex =
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    if (!userForm.email.trim())
      errors.email = "Please, enter your email.";
    else if (!emailRegex.test(userForm.email))
      errors.email = "Please, enter a valid email.";

    if (!userForm.password)
      errors.password = "Please, enter your password.";

    if (Object.keys(errors).length > 0) {
      setErrors(errors);
      return;
    }

    const res = await fetch(`${BASE_API_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify(userForm)
    });

    if (!res.ok) {
      switch (res.status) {
        case 500: {
          errors.req = "There was an error on the server. Please, try again later."
        } break;
        case 401: {
          errors.req = "Invalid combination of email and password. Double-check both fields and try again."
        } break;
        default: {
          errors.req = "There was an error on the server. Please, try again later."
        }
      }
      setErrors(errors)
      return;
    }

    setMe((await res.json())["user"]);
    navigate("/main");
  };

  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <div className="auth-two-col">
          <div className="auth-two-col-title">
            <h1>Welcome back! Hope you're about to have a good playing session!</h1>
          </div>
          <div className="auth-two-col-form">
            <h2>Log in</h2>
            { errors.req && <div className="form-error">{errors.req}</div> }
            <form onSubmit={handleSubmit} noValidate>
              <div className="form-entry">
                <label htmlFor="email">Email</label>
                <input required id="email" name="email" type="email" />
                { errors.email && <div className="form-entry-error">{errors.email}</div> }
              </div>
              <div className="form-entry">
                <label htmlFor="password">Password</label>
                <input required id="password" name="password" type="password" />
                { errors.password && <div className="form-entry-error">{errors.password}</div> }
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

