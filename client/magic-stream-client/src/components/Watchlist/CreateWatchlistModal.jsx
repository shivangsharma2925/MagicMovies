import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faX,
  faBookmark,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";
import { useCreateWatchlist } from "../../hook/UseWatchlists";
import "../../styling/WatchList.css";

export default function CreateWatchlistModal({ onClose }) {
  const [name, setName] = useState("");
  const [formError, setFormError] = useState("");

  const createMutation = useCreateWatchlist();

  const handleCreate = async () => {
    const watchlistName = name.trim();

    if (!watchlistName) {
      setFormError("Please enter a watchlist name.");
      return;
    }

    await createMutation.mutateAsync(watchlistName);

    setName("");
    setFormError("");
    onClose?.();
  };

  return (
    <div
      className="wl-modal-bg"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose?.();
      }}
    >
      <div className="wl-modal wl-create-modal">
        <div className="wl-modal-header">
          <span className="wl-modal-title">Create Watchlist</span>

          <button
            type="button"
            className="wl-close-btn"
            onClick={() => onClose?.()}
            aria-label="Close"
          >
            <FontAwesomeIcon icon={faX} />
          </button>
        </div>

        <div className="wl-create-hero">
          <div className="wl-empty-icon-ring small">
            <FontAwesomeIcon icon={faBookmark} />
          </div>

          <h2>Create your first list</h2>
          <p>
            Organize movies into personal lists and track what you want to watch.
          </p>
        </div>

        <div className="wl-create-form-card">
          <label className="wl-create-label">Watchlist name</label>

          <input
            autoFocus
            className="wl-create-input"
            placeholder="Example: Must Watch"
            value={name}
            onChange={(e) => {
              setName(e.target.value);
              setFormError("");
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter") handleCreate();
              if (e.key === "Escape") onClose?.();
            }}
          />

          {formError && <p className="wl-form-error">{formError}</p>}

          <button
            type="button"
            className="wl-create-submit"
            onClick={handleCreate}
            disabled={createMutation.isPending}
          >
            <FontAwesomeIcon icon={faPlus} />
            {createMutation.isPending ? "Creating..." : "Create Watchlist"}
          </button>
        </div>
      </div>
    </div>
  );
}
