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
  }, []);

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
      <Form className="mb-3">
        <Form.Control
          type="text"
          placeholder="Search in Recommended movies..."
          value={search}
          className="w-50 m-auto mt-2"
          onChange={(e) => setSearch(e.target.value)}
        />
      </Form>

      {loading ? (
        <h2><Spinner /></h2>
      ) : (
        <Movies movies={filteredMovies} message={message} />
      )}
    </>
  );
};

export default Recommended;