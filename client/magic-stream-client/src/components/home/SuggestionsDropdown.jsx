import { ListGroup, Image } from "react-bootstrap";

const SuggestionsDropdown = ({ suggestions, onSuggestionClick }) => {
  if (!suggestions?.length) return null;

  return (
    <div className="suggestions-dropdown">
      {suggestions.map((movie) => (
        <div
          key={movie.imdb_id}
          action
          onClick={() => onSuggestionClick(movie.title)}
          className="suggestion-item"
        >
          <img
            src={movie.poster_path}
            alt={movie.title}
            className="suggestion-poster"
          />

          <span className="suggestion-title">{movie.title}</span>
        </div>
      ))}
    </div>
  );
};

export default SuggestionsDropdown;
