import { useEffect, useState, useMemo } from "react";
import { Form } from "react-bootstrap";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import Movies from "../movies/Movies";
import Spinner from "../../utils/Spinner";

const Recommended = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState("");

  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchRecommendedMovies = async () => {
      setLoading(true);
      try {
        const resp = await axiosPrivate.get("/recommendedmovies");
        setMovies(resp.data);
      } catch (error) {
        console.log(error);
      } finally {
        setLoading(false);
      }
    };

    fetchRecommendedMovies();
  }, [axiosPrivate]);

  // derived state
  const filteredMovies = useMemo(() => {
    if (!search.trim()) return movies;

    return movies.filter(movie =>
      movie.title.toLowerCase().includes(search.toLowerCase())
    );
  }, [movies, search]);

  // derived message
  const message = useMemo(() => {
    if (loading) return "";
    if (!movies.length) return "No Movies found!!";
    if (search && !filteredMovies.length)
      return "No Movies found, Please search again!";
    return "";
  }, [movies, filteredMovies, search, loading]);

  return (
    <>
      <div className="recommended-search-wrap">
        <span className="search-field-icon" aria-hidden="true">&#128269;</span>
          <Form.Control
            type="text"
            placeholder="Search in Recommended movies..."
            value={search}
            className="search-field"
            onChange={(e) => setSearch(e.target.value)}
          />
      </div>

      {loading ? (
        <h2><Spinner /></h2>
      ) : (
        <Movies movies={filteredMovies} message={message} />
      )}
    </>
  );
};

export default Recommended;