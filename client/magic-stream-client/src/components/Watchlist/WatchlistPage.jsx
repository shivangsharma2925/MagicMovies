import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBookmark,
  faPlus,
  faList,
  faFilm,
  faEye,
} from "@fortawesome/free-solid-svg-icons";
import { useWatchlists } from "../../hook/UseWatchlists";
import CreateWatchlistModal from "./CreateWatchlistModal";
import WatchlistPanel from "./WatchListPanel";
import PageTitle from "../../utils/PageTitle";

export default function WatchlistPage() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [autoModalClosed, setAutoModalClosed] = useState(false);

  const { data: watchLists = [], isLoading } = useWatchlists();

  const isEmpty = !isLoading && watchLists.length === 0;

  const shouldShowCreateModal =
    isCreateModalOpen || (isEmpty && !autoModalClosed);

  const closeCreateModal = () => {
    setIsCreateModalOpen(false);
    setAutoModalClosed(true);
  };

  if (isLoading) {
    return <p className="wl-loading">Loading watchlists...</p>;
  }

  return (
    <>
      <PageTitle title="Watchlists" />
      {shouldShowCreateModal && (
        <CreateWatchlistModal onClose={closeCreateModal} />
      )}

      {isEmpty && (
        <div className="wl-empty-page">
          <div className="wl-empty-card">
            <div className="wl-empty-icon-ring">
              <FontAwesomeIcon icon={faBookmark} />
            </div>

            <h2 className="wl-empty-title">No watchlists yet</h2>

            <p className="wl-empty-desc">
              Create your first watchlist and start saving movies you want to
              watch later.
            </p>

            <div className="wl-steps">
              <div className="wl-step">
                <div className="wl-step-num">1</div>
                <FontAwesomeIcon icon={faList} className="wl-step-icon" />
                <span className="wl-step-label">Create a list</span>
              </div>

              <div className="wl-step">
                <div className="wl-step-num">2</div>
                <FontAwesomeIcon icon={faFilm} className="wl-step-icon" />
                <span className="wl-step-label">Add movies</span>
              </div>

              <div className="wl-step">
                <div className="wl-step-num">3</div>
                <FontAwesomeIcon icon={faEye} className="wl-step-icon" />
                <span className="wl-step-label">Track watched</span>
              </div>
            </div>

            <button
              type="button"
              className="wl-empty-create-btn"
              onClick={() => setIsCreateModalOpen(true)}
            >
              <FontAwesomeIcon icon={faPlus} />
              Create Watchlist
            </button>
          </div>
        </div>
      )}

      {!isEmpty && <WatchlistPanel />}
    </>
  );
}