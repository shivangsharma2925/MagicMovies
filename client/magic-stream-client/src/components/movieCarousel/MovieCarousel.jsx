import { useRef } from "react";
import { useNavigate } from "react-router-dom";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlay, faCircleInfo } from "@fortawesome/free-solid-svg-icons";
import { slugify } from "../../utils/SlugifyTitle";
import "../../styling/Carousel.css";

export default function MovieCarousel({ title, movies, section }) {
  const trackRef = useRef(null);
  const navigate = useNavigate();

  const scroll = (dir) => {
    trackRef.current?.scrollBy({ left: dir * 420, behavior: "smooth" });
  };

  if (!movies?.length) return null;

  return (
    <div className="carousel-section">
      <div className="carousel-header">
        <h3 className="carousel-title">{title}</h3>
        <span className="carousel-see-all" onClick={() => navigate(`/category/${section}`)}>See all →</span>
      </div>

      <div className="carousel-track-wrap">
        <button
          className="carousel-arrow arrow-left"
          onClick={() => scroll(-1)}
        >
          ‹
        </button>

        <div className="carousel-track" ref={trackRef}>
          {movies.map((movie) => (
            <div
              key={movie._id}
              className="mc-card"
              onClick={() => navigate(`/stream/${slugify(movie.title)}/${movie.imdb_id}`)}
            >
              <div className="mc-poster">
                <img src={movie.poster_path} alt={movie.title} loading="lazy" />
                <div className="mc-overlay">
                  <div className="mc-play-btn">
                    <FontAwesomeIcon icon={faPlay} />
                  </div>
                </div>
                {movie.ranking?.ranking_name && (
                  <span className="mc-rank-badge">
                    {movie.ranking.ranking_name}
                  </span>
                )}
              </div>
              <div className="mc-info">
                <div className="mc-title">{movie.title}</div>
                <div className="mc-meta">
                  {movie.release_date?.split("-")[0]}{" "}
                  {/* extract year from "1994-09-23" */}
                  {movie.runtime &&
                    ` · ${Math.floor(movie.runtime / 60)}h ${movie.runtime % 60}m`}
                  <button
                    className="mc-info-icon"
                    onClick={(e) => {e.stopPropagation(); navigate(`/movie/detail/${movie.imdb_id}`)}}
                    aria-label="More info"
                    title="More Info"
                  >
                    <FontAwesomeIcon icon={faCircleInfo} />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        <button
          className="carousel-arrow arrow-right"
          onClick={() => scroll(1)}
        >
          ›
        </button>
      </div>
    </div>
  );
}
