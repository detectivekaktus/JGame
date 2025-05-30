import { Spinner } from "../components/Spinner";

export function LoadingPage() {
  return (
    <div className="page-wrapper">
      <div className="page content">
        <Spinner />
      </div>
    </div>
  );
}
