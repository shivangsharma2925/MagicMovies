import { useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlay} from "@fortawesome/free-solid-svg-icons";

export default function SearchedMovies({ movies = [], search, setSearchOpen, setSearch, setShowSuggestions }) {
  const navigate = useNavigate();

  if (!search || search.trim().length < 2) return null;

  return (
    <div className="sm-wrap">
      {movies && movies.length > 0 ? (
        <>
          <div className="sm-divider">
            <span>{movies.length} result{movies.length !== 1 ? "s" : ""} for "{search}"</span>
            <hr />
          </div>

          <div className="sm-track">
            {movies.map((movie) => (
              <div
                key={movie._id}
                className="sm-card"
                onClick={() => {navigate(`/stream/${movie.youtube_id}`); setSearchOpen(false); setSearch(""); setShowSuggestions(false);}}
              >
                <div className="sm-poster">
                  <img src={movie.poster_path} alt={movie.title} loading="lazy" />
                  <div className="sm-overlay">
                    <div className="sm-play">
                      <FontAwesomeIcon icon={faPlay} />
                    </div>
                  </div>
                </div>
                <div className="sm-info">
                  <div className="sm-title">{movie.title}</div>
                  <div className="sm-year">{movie.release_date?.split("-")[0]}</div>
                </div>
              </div>
            ))}
          </div>
        </>
      ) : (
        <div className="sm-empty">
          Uh oh, No movies found for "{search}" !!
        </div>
      )}
    </div>
  );
}