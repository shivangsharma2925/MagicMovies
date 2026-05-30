import { ListGroup, Image } from "react-bootstrap";

const SuggestionsDropdown = ({ suggestions, onSuggestionClick }) => {
  if (!suggestions?.length) return null;

  return (
    <div
      className="position-absolute start-50 translate-middle-x mt-2"
      style={{
        width: "100%",
        zIndex: 3000,
      }}
    >
      <ListGroup
        className="shadow border-0 overflow-hidden suggestionsList"
      >
        {suggestions.map((movie) => (
          <ListGroup.Item
            key={movie.imdb_id}
            action
            onClick={() => onSuggestionClick(movie.title)}
            className="d-flex align-items-center gap-3 px-3 py-2 border-0 suggestion-item"
            style={{
              transition: "all 0.2s ease",
              cursor: "pointer",
            }}
          >
            <Image
              src={movie.poster_path}
              width={50}
              height={50}
              rounded
              style={{
                objectFit: "cover",
                flexShrink: 0,
              }}
            />

            <div
              className="fw-medium text-truncate"
              style={{
                fontSize: "1rem",
              }}
            >
              {movie.title}
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
};

export default SuggestionsDropdown;
