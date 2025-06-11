import { Button } from "./Button";
import "../css/Search.css"
import React from "react";

type SearchProps = {
  placeholder: string,
  setQuery: React.Dispatch<React.SetStateAction<string>>
  handleQuery: () => void
}

export function Search({ placeholder, setQuery, handleQuery }: SearchProps) {
  return (
    <div className="margin-top container search">
      <input className="search-bar" onChange={(e) => setQuery(e.currentTarget.value)} name="query" id="query" type="text" placeholder={placeholder} />
      <Button onClick={handleQuery} className="search-bar-submit" stretch={false} dim={false}>Search</Button>
    </div>
  );
}
