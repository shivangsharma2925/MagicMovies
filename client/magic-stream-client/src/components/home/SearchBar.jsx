import SuggestionsDropdown from "./SuggestionsDropdown";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMagnifyingGlass, faXmark } from "@fortawesome/free-solid-svg-icons";

const SearchBar = ({
  search,
  setSearch,
  suggestions,
  showSuggestions,
  setShowSuggestions,
}) => {

  return (
    <div className={`nb-search-box ${showSuggestions ? "search-input-wrapper" : ""}`}>
      <FontAwesomeIcon icon={faMagnifyingGlass} className="nb-search-icon" />
      <input
        className="search-field"
        type="text"
        placeholder="Search movies..."
        value={search}
        onFocus={() => {
          if (suggestions.length > 0) {
            setShowSuggestions(true);
          }
        }}
        onChange={(e) => {
          setSearch(e.target.value);

          if (e.target.value.trim().length >= 2) {
            setShowSuggestions(true);
          } else {
            setShowSuggestions(false);
          }
        }}
      />

      <button
        className="nb-search-clear"
        onClick={()=>{setSearch(""); setShowSuggestions(false);}}
        aria-label="clear search"
      >
        <FontAwesomeIcon icon={faXmark} />
      </button>
    </div>
  );
};

export default SearchBar;
