import { useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faXmark, faCheck, faFilm, faPlus } from "@fortawesome/free-solid-svg-icons";
import {useWatchlists, useAddToWatchlist, useRemoveFromWatchlist, } from "../../hook/UseWatchlists";
import "../../styling/WatchList.css";

const COLORS = [
  { bg: "rgba(59,130,246,0.18)", text: "#60a5fa" },
  { bg: "rgba(34,197,94,0.18)", text: "#4ade80" },
  { bg: "rgba(234,179,8,0.18)", text: "#fbbf24" },
  { bg: "rgba(168,85,247,0.18)", text: "#c084fc" },
  { bg: "rgba(239,68,68,0.18)", text: "#f87171" },
];

export default function AddToWatchlistPicker({ movie, onClose }) {
  const navigate = useNavigate();
  const { data: watchLists = [], isLoading } = useWatchlists();
  const addMutation = useAddToWatchlist();
  const removeMutation = useRemoveFromWatchlist();

  const handleToggle = (list) => {
    const exists = list.movies?.some(m => m.imdb_id === movie.imdb_id);
    if (exists) {
      removeMutation.mutate({ listId: list.id, imdbId: movie.imdb_id });
    } else {
      addMutation.mutate({
        listId: list.id,
        movie: {
          imdb_id: movie.imdb_id,
          title:   movie.title,
          poster:  movie.poster_path,
          youtube_id : movie.youtube_id,
          year:    movie.release_date?.split("-")[0] ?? "",
          watched: false,
        },
      });
    }
  };

  const goCreateList = () => {
    onClose();
    navigate("/watchlist"); // full page handles creation
  };

  return (
    <div className="atw-overlay" onClick={e => e.target === e.currentTarget && onClose()}>
      <div className="atw-modal">

        <div className="atw-header">
          <span className="atw-title">Add to Watchlist</span>
          <button className="atw-close" onClick={onClose} aria-label="Close">
            <FontAwesomeIcon icon={faXmark} />
          </button>
        </div>
        <p className="atw-movie-name">"{movie?.title}"</p>

        {isLoading ? (
          <p className="atw-loading">Loading...</p>
        ) : watchLists.length === 0 ? (
          // No lists exist — redirect to full watchlist page
          <div className="atw-empty">
            <div className="atw-empty-icon">
              <FontAwesomeIcon icon={faFilm} />
            </div>
            <p>You don't have any watchlists yet</p>
            <span>Create your first list to start saving movies</span>
            <button className="atw-goto-btn" onClick={goCreateList}>
              <FontAwesomeIcon icon={faPlus} /> Create a Watchlist
            </button>
          </div>
        ) : (
          <div className="atw-list">
            {watchLists.map((list, i) => {
              const c      = COLORS[i % COLORS.length];
              const exists = list.movies?.some(m => m.imdb_id === movie?.imdb_id);
              const isAddingThisList = addMutation.isPending && addMutation.variables?.listId === list.id;
              const isRemovingThisList = removeMutation.isPending && removeMutation.variables?.listId === list.id;
              const isUpdatingThisList = isAddingThisList || isRemovingThisList;

              return (
                <button
                  disabled={isUpdatingThisList}
                  key={list.id}
                  className={`atw-item ${exists ? "checked" : ""} ${isUpdatingThisList ? "disabled" : ""}`}
                  onClick={() => handleToggle(list)}
                >
                  <div className="atw-item-icon" style={{ background: c.bg, color: c.text }}>
                    <FontAwesomeIcon icon={faFilm} />
                  </div>
                  <span className="atw-item-name">{list.name}</span>
                  <span className="atw-item-count">
                    {list.movies?.length ?? 0} film{list.movies?.length !== 1 ? "s" : ""}
                  </span>
                  <div className={`atw-checkbox ${exists ? "checked" : ""}`}>
                    {exists && <FontAwesomeIcon icon={faCheck} />}
                  </div>
                </button>
              );
            })}
          </div>
        )}

        {watchLists.length > 0 && (
          <button className="atw-manage-link" onClick={goCreateList}>
            Manage all watchlists →
          </button>
        )}

      </div>
    </div>
  );
}