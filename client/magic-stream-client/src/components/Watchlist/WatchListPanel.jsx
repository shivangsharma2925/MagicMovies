import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faBookmark,
  faPen,
  faTrash,
  faEye,
  faEyeSlash,
  faCircleMinus,
  faChevronDown,
  faCheck,
  faFilm,
} from "@fortawesome/free-solid-svg-icons";
import {
  useWatchlists,
  useCreateWatchlist,
  useDeleteWatchlist,
  useRenameWatchlist,
  // useAddToWatchlist,
  useRemoveFromWatchlist,
  useToggleWatched,
} from "../../hook/UseWatchlists";
import "../../styling/WatchList.css";

const COLORS = [
  { bg: "rgba(59,130,246,0.18)", text: "#60a5fa" },
  { bg: "rgba(34,197,94,0.18)", text: "#4ade80" },
  { bg: "rgba(234,179,8,0.18)", text: "#fbbf24" },
  { bg: "rgba(168,85,247,0.18)", text: "#c084fc" },
  { bg: "rgba(239,68,68,0.18)", text: "#f87171" },
];

export default function WatchlistPanel({ movie }) {
  const [openIds, setOpenIds] = useState([]);
  const [creating, setCreating] = useState(false);
  const [newName, setNewName] = useState("");
  const [renamingId, setRenamingId] = useState(null);
  const [renameName, setRenameName] = useState("");
  const [deletingId, setDeletingId] = useState(null);
  const [actionError, setActionError] = useState("");

  const navigate = useNavigate();

  const { data: watchLists = [], isLoading } = useWatchlists();

  const createMutation = useCreateWatchlist();
  const deleteMutation = useDeleteWatchlist();
  const renameMutation = useRenameWatchlist();
  // const addMutation = useAddToWatchlist();
  const removeMutation = useRemoveFromWatchlist();
  const watchedMutation = useToggleWatched();

  const toggleOpen = (id) =>
    setOpenIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id],
    );

  const handleCreate = async () => {
    const trimmedName = newName.trim();

    if (!trimmedName) {
      setActionError("Please enter a watchlist name.");
      return;
    }

    try {
      setActionError("");

      await createMutation.mutateAsync(trimmedName);
    } catch (error) {
      const message =
        error.response?.data?.error ||
        "Unable to create watchlist. Please try again.";

      setActionError(message);
    } finally {
      setNewName("");
      setCreating(false);
    }
  };

  const handleRename = async (id) => {
    if (!renameName.trim()) {
      setActionError("List name cannot be empty.");
      return;
    }
    await renameMutation.mutateAsync({ id: id, name: renameName.trim() });
    setRenamingId(null);
    setRenameName("");
  };

  const handleDelete = (id) => setDeletingId(id);
  const confirmDelete = () => {
    deleteMutation.mutate(deletingId);
    setDeletingId(null);
  };

  // const handleAdd = (listId) => {
  //   const list = watchLists.find((l) => l.id === listId);
  //   if (list?.movies?.some((m) => m.imdb_id === movie?.imdb_id)) return;
  //   addMutation.mutate({
  //     watchListId: listId,
  //     movie: {
  //       imdb_id: movie.imdb_id,
  //       title: movie.title,
  //       poster: movie.poster_path,
  //       youtube_id : movie.youtube_id,
  //       year: movie.release_date?.split("-")[0] ?? "",
  //       watched: false,
  //     },
  //   });
  // };

  const handleRemove = (listId, movieId) =>
    removeMutation.mutate({ listId: listId, imdbId: movieId });

  const handleToggleWatched = (listId, movieId) =>
    watchedMutation.mutate({
      listId: listId,
      imdbId: movieId,
    });

  if (isLoading)
    return <div className="wlp-loading">Loading watchlists...</div>;

  return (
    <div className="wlp-wrap">
      {/* Header row */}
      <div className="wlp-header">
        <div className="wlp-header-left">
          <FontAwesomeIcon icon={faBookmark} className="wlp-header-icon" />
          <span className="wlp-header-title">Watchlists</span>
          {movie && (
            <span className="wlp-header-movie">for "{movie.title}"</span>
          )}
        </div>
        <button className="wlp-new-btn" onClick={() => setCreating(true)}>
          <FontAwesomeIcon icon={faPlus} /> New list
        </button>
      </div>

      {/* Create new input */}
      {creating && (
        <div className="wlp-action-error-wrapper">
          <div className="wlp-create-row">
            <input
              autoFocus
              className="wlp-create-input"
              placeholder="List name..."
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleCreate();
                if (e.key === "Escape") {
                  setCreating(false);
                  setNewName("");
                }
              }}
            />
            <button
              className="wlp-btn-sm accent"
              onClick={handleCreate}
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? "..." : "Create"}
            </button>
            <button
              className="wlp-btn-sm"
              onClick={() => {
                setCreating(false);
                setNewName("");
              }}
            >
              Cancel
            </button>
          </div>
          {actionError && (
            <p
              className="wl-form-error"
              style={{ padding: "0px 0px 0px 17px" }}
            >
              {actionError}
            </p>
          )}
        </div>
      )}

      {/* Delete confirm banner */}
      {deletingId && (
        <div className="wlp-delete-confirm">
          <span>
            Delete{" "}
            <b className="text-white">
              "{watchLists.find((l) => l.id === deletingId)?.name}"
            </b>{" "}
            ? This can't be undone.
          </span>
          <div className="wlp-delete-btns">
            <button className="wlp-btn-sm" onClick={() => setDeletingId(null)}>
              Cancel
            </button>
            <button className="wlp-btn-sm danger" onClick={confirmDelete}>
              Delete
            </button>
          </div>
        </div>
      )}

      {/* Accordion list */}
      <div className="wlp-accordions">
        {watchLists.map((list, i) => {
          const c = COLORS[i % COLORS.length];
          const isOpen = openIds.includes(list.id);
          const entry = list.movies?.find((m) => m.imdb_id === movie?.imdb_id);
          const exists = !!entry;
          const isRenaming = renamingId === list.id;

          return (
            <div
              key={list.id}
              className={`wlp-accordion ${exists ? "has-movie" : ""}`}
              style={{ background: c.bg }}
            >
              {/* Accordion header */}
              <div className="wlp-acc-header">
                {isRenaming ? (
                  <div className="wlp-action-error-wrapper">
                    <div className="wlp-rename-row">
                      <input
                        autoFocus
                        className="wlp-rename-input"
                        value={renameName}
                        onChange={(e) => setRenameName(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") handleRename(list.id);
                          if (e.key === "Escape") {
                            setRenamingId(null);
                            setRenameName("");
                            setActionError("");
                          }
                        }}
                        placeholder="Enter New Name..."
                      />
                      <button
                        className="wlp-btn-sm accent"
                        onClick={() => handleRename(list.id)}
                      >
                        Save
                      </button>
                      <button
                        className="wlp-btn-sm"
                        onClick={() => {
                          setRenamingId(null);
                          setRenameName("");
                          setActionError("");
                        }}
                      >
                        Cancel
                      </button>
                    </div>
                    {actionError && (
                      <p
                        className="wl-form-error"
                        style={{ padding: "0px 0px 0px 17px" }}
                      >
                        {actionError}
                      </p>
                    )}
                  </div>
                ) : (
                  <>
                    <button
                      className="wlp-acc-toggle"
                      onClick={() => toggleOpen(list.id)}
                    >
                      <div
                        className="wlp-acc-icon"
                        style={{ background: c.bg, color: c.text }}
                      >
                        <FontAwesomeIcon icon={faFilm} />
                      </div>
                      <div className="wlp-acc-info">
                        <span
                          className="wlp-acc-name"
                          style={{ color: c.text }}
                        >
                          {list.name}
                        </span>
                        <span className="wlp-acc-meta">
                          {list.movies?.length ?? 0} film
                          {list.movies?.length !== 1 ? "s" : ""}
                          {exists && <span className="wlp-in-list-dot" />}
                          {exists && (
                            <span style={{ color: c.text }}>In this list</span>
                          )}
                        </span>
                      </div>
                      <FontAwesomeIcon
                        icon={faChevronDown}
                        className={`wlp-chevron ${isOpen ? "open" : ""}`}
                      />
                    </button>

                    <div className="wlp-acc-actions">
                      {/* Rename List */}
                      <button
                        className="wlp-icon-btn"
                        onClick={() => {
                          setRenamingId(list.id);
                          setRenameName(list.name);
                          setDeletingId(null);
                          setActionError("");
                        }}
                        title="Rename list"
                        aria-label="Rename"
                      >
                        <FontAwesomeIcon icon={faPen} />
                      </button>

                      {/* Delete List */}
                      <button
                        className="wlp-icon-btn danger"
                        onClick={() => handleDelete(list.id)}
                        title="Delete list"
                        aria-label="Delete"
                      >
                        <FontAwesomeIcon icon={faTrash} />
                      </button>
                    </div>
                  </>
                )}
              </div>

              {/* Accordion body — movie list */}
              <div className={`wlp-acc-panel ${isOpen ? "open" : ""}`}>
                <div className="wlp-acc-panel-inner">
                  <div className="wlp-acc-body">
                    {list.movies?.length === 0 ? (
                      <div className="wlp-acc-empty">
                        No movies in this list yet
                      </div>
                    ) : (
                      list.movies.map((m) => (
                        <div
                          key={m.imdb_id}
                          className={`wlp-movie-row ${m.imdb_id === movie?.imdb_id ? "current" : ""}`}
                        >
                          <div
                            className="wlp-movie-thumb"
                            onClick={() => navigate(`/stream/${m.youtube_id}`)}
                          >
                            {m.poster ? (
                              <img src={m.poster} alt={m.title} />
                            ) : (
                              <FontAwesomeIcon icon={faFilm} />
                            )}
                          </div>

                          <div className="wlp-movie-info">
                            <span className="wlp-movie-title">{m.title}</span>
                            <span className="wlp-movie-meta">
                              {m.year}
                              {m.watched && (
                                <span className="wlp-watched-tag">
                                  <FontAwesomeIcon icon={faCheck} /> Watched
                                </span>
                              )}
                            </span>
                          </div>
                          <div className="wlp-acc-actions">
                            <button
                              className={`wlp-icon-btn ${m?.watched ? "watched" : ""}`}
                              onClick={() =>
                                handleToggleWatched(list.id, m.imdb_id)
                              }
                              aria-label={
                                m?.watched
                                  ? "Mark as unwatched"
                                  : "Mark as watched"
                              }
                              title={
                                m?.watched
                                  ? "Mark as unwatched"
                                  : "Mark as watched"
                              }
                            >
                              {m?.watched ? (
                                <FontAwesomeIcon icon={faEye} />
                              ) : (
                                <FontAwesomeIcon icon={faEyeSlash} />
                              )}
                            </button>

                            <button
                              className="wlp-icon-btn warning"
                              onClick={() => handleRemove(list.id, m.imdb_id)}
                              aria-label="Remove from List"
                              title="Remove from List"
                            >
                              <FontAwesomeIcon icon={faCircleMinus} />
                            </button>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
