import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { MeContext } from "../context/MeProvider";
import { LoadingPage } from "./LoadingPage";
import { Button } from "../components/Button";
import { Footer } from "../components/Footer";
import { BASE_API_URL } from "../utils/consts";
import "../css/SettingsPage.css"
import "../css/Form.css"

type SettingsFormType = {
  name: string
  email: string
  newPassword: string
  confirmNewPassword: string
}

type PatchRequstType = {
  name?: string,
  email?: string,
  password?: string
}

export function SettingsPage() {
  const { me, setMe, loadingMe } = useContext(MeContext);
  const navigate = useNavigate();
  const [form, setForm] = useState<SettingsFormType | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (loadingMe)
      return;

    if (!me) {
      navigate("/auth/login");
      return;
    }

    setForm({
      name: me.name,
      email: me.email,
      newPassword: "",
      confirmNewPassword: ""
    });
  }, [me, loadingMe]);

  if (loadingMe)
    return <LoadingPage />

  const handleChange = (field: string) => 
    (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((prev) => ({ ...prev as SettingsFormType, [field]: e.target.value }));

  const handleSubmit = async () => {
    if (!form)
      throw new Error("Form is null. Something went wrong");
    else if (!me)
      throw new Error("Me is null. Something went wrong.")

    const req = { } as PatchRequstType;
    const newErrors: Record<string, string> = {}

    const name = form.name.trim();
    if (name !== me.name) {
      if (name.length < 4 || name.length > 32) {
        newErrors.name = "Name must be between 4 and 32 characters.";
        return;
      }
      req.name = name;
    }

    const email = form.email.trim();
    if (email !== me.email) {
      const emailRegex =
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
      if (!emailRegex.test(email)) {
        newErrors.email = "Please, enter a valid email.";
        return;
      }
      req.email = email;
    }

    const newPassword = form.newPassword;
    const confirmNewPassword = form.confirmNewPassword;
    if (newPassword.length !== 0 && confirmNewPassword.length !== 0) {
      if (newPassword !== confirmNewPassword) {
        newErrors.password = "Passwords don't match."
        return;
      }
      req.password = newPassword;
    }

    if (Object.keys(req).length === 0) {
      newErrors.req = "No changes detected."
      setErrors(newErrors);
      return;
    }

    const res = await fetch(`${BASE_API_URL}/users/me`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify(req)
    })
    const body = await res.json();

    if (!res.ok) {
      newErrors.req = body["message"];
      setErrors(newErrors);
      return;
    }

    setMe(body);
    setErrors({});
    document.location.reload();
  }

  const handleDelete = async () => {
    const newErrors: Record<string, string> = {}

    const res = await fetch(`${BASE_API_URL}/users/me`, {
      method: "DELETE",
      credentials: "include"
    })
    const body = await res.json();

    if (!res.ok) {
      errors.req = body["message"];
      setErrors(newErrors);
      return;
    }

    setMe(null);
    setErrors({});
    document.location.reload();
  };

  return (
    <div className="page-wrapper">
      <div className="margin-top margin-bottom container page">
        <h1>Settings</h1>
        { errors.req && <div className="form-error">{errors.req}</div> }
        <hr />
        <dl>
          <div className="dl-section">
            <dt>Public information</dt>
            <dd>
              <p>Your name</p>
              <input name="name" type="text" defaultValue={me?.name} onChange={handleChange("name")}/>
            </dd>
          </div>
          <div className="dl-section">
            <dt>Private information</dt>
            <dd>
              <p>Your email</p>
              <input name="email" type="email" defaultValue={me?.email} onChange={handleChange("email")}/>
            </dd>
          </div>
          <div className="dl-section">
            <dt>Password</dt>
            <dd>
              <p>New password</p>
              <input name="newPassword" type="password" onChange={handleChange("newPassword")}/>
              <p>Re-type new password</p>
              <input name="confirmNewPassword" type="password" onChange={handleChange("confirmNewPassword")}/>
            </dd>
          </div>
        </dl>
        <div className="button-options">
          <Button stretch={true} dim={false} onClick={() => navigate(-1)}>← Go back</Button>
          <Button className="save-changes-button" stretch={true} dim={false} onClick={handleSubmit}>Save changes</Button>
        </div>
        <hr />
        <Button stretch={true} dim={false} onClick={handleDelete}>Delete account</Button>
      </div>
      <Footer />
    </div>
  );
}
