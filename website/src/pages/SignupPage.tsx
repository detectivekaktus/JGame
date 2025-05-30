import React, { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Footer } from "../components/Footer"
import { SignupForm } from "../types/user";
import { BASE_API_URL } from "../utils/consts";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import "../css/Form.css"

export function SignupPage() {
  const { me, loadingMe } = useContext(MeContext);

  const [errors, setErrors] = useState<Record<string, string>>({});
  const navigate = useNavigate();

  useEffect(() => {
    if (loadingMe)
      return;

    if (me)
      navigate("/main");
  }, [me, loadingMe]);

  if (loadingMe)
    return <LoadingPage />

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const target: HTMLFormElement = e.currentTarget;
    const data = new FormData(target);
    const userForm: SignupForm = {
      name: data.get("name") as string,
      email: data.get("email") as string,
      password: data.get("password") as string,
      confirm_password: data.get("confirm-password") as string
    };
    const newErrors: Record<string, string> = {}

    const username = userForm.name.trim();
    if (!username)
      newErrors.name = "Please, enter your name.";
    else if (username.length < 4 || username.length > 32)
      newErrors.name = "Name must be between 4 and 32 characters.";

    // Copied from StackOverflow
    const emailRegex =
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    const email = userForm.email.trim()
    if (!email)
      newErrors.email = "Please, enter your email.";
    else if (!emailRegex.test(email))
      newErrors.email = "Please, enter a valid email.";

    if (!userForm.password)
      newErrors.password = "Please, enter your password.";
    else if (userForm.password.length < 8 || userForm.password.length > 32)
      newErrors.password = "Password must be between 8 and 32 characters long.";
    else if (!/^(?=.*[a-z])+(?=.*[A-Z])+(?=.*[0-9])+(?=.*[\!\$\@\#\^\&]).{8,32}$/.test(userForm.password))
      newErrors.password = "Password must contain uppercase, lowercase, number and one special symbol.";

    if (!userForm.confirm_password)
      newErrors.confirm_password = "Please, repeat your password.";
    else if (userForm.password !== userForm.confirm_password)
      newErrors.confirm_password = "Passwords do not match.";

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    const res = await fetch(`${BASE_API_URL}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name: userForm.name,
        email: userForm.email,
        password: userForm.password
      })
    });
    const body = await res.json();

    if (!res.ok) {
      newErrors.req = body["message"];
      setErrors(newErrors);
      return;
    }
    else
      navigate("/auth/login");
  };

  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <div className="auth-two-col">
          <div className="auth-two-col-title">
            <h1>Good to see you! Let’s sign up, it’s quick and simple</h1>
          </div>
          <div className="auth-two-col-form">
            <h2>Sign up</h2>
            { errors.req && <div className="form-error">{errors.req}</div> }
            <form onSubmit={handleSubmit} noValidate>
              <div className="form-entry">
                <label htmlFor="name">Display name</label>
                <input required id="name" name="name" type="text" />
                { errors.name && <div className="form-entry-error">{errors.name}</div> }
              </div>
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
              <div className="form-entry">
                <label htmlFor="confirm-password">Repeat password</label>
                <input required id="confirm-password" name="confirm-password" type="password" />
                { errors.confirm_password && <div className="form-entry-error">{errors.confirm_password}</div> }
              </div>
              <Button stretch={true} dim={false} type="submit">Sign up</Button>
            </form>
            <p>Already have an account?  <Link className="fg-accent-600 underlined" to={"/auth/login"}>Log in</Link></p>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
}

