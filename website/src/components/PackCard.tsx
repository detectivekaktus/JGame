import { Button } from "./Button";
import { useEffect, useState } from "react";
import "../css/PackCard.css"
import { BASE_API_URL } from "../utils/consts";

export type Pack = {
  id: number,
  name: string,
  user_id: number,
  body: object
}

export type PackCardProps = {
  pack: Pack
}

export function PackCard({ pack }: PackCardProps) {
  const [popupDisplayed, setPopupDisplayed] = useState(false);
  const [packOwnerName, setPackOwnerName] = useState("unknown");

  useEffect(() => {
    fetch(`${BASE_API_URL}/users/${pack.user_id}`)
      .then(res => res.status === 200 ? res.json() : null)
      .then(data => setPackOwnerName(data["name"]))
      .catch(err => console.error(err))
  }, [])

  const showPopup = () => setPopupDisplayed(true);
  const hidePopup = () => {
    setPopupDisplayed(false);
    navigator.clipboard.writeText(String(pack.id))
      .catch(err => console.error(err))
  }

  return (
    <>
      <Button onClick={showPopup} className="pack-card" stretch={true} dim={false}>
        <strong><span className="pack-card-name">{pack.name}</span></strong>
        <span className="pack-card-owner">Created by: {packOwnerName}</span>
      </Button>
      {
        popupDisplayed && (
          <div className="pack-card-popup">
            <div className="pack-card-popup-content">
              <h3>Pack id: {pack.id}</h3>
              <Button onClick={hidePopup} stretch={false} dim={false}>Close and copy id</Button>
            </div>
          </div>
        )
      }
    </>
  );
}
