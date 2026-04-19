import { useEffect, useState } from "react";
import { Form } from "react-bootstrap";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import Movies from "../movies/Movies";
import Spinner from "../../utils/Spinner";

const Recommended = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();
    const [search, setSearch] = useState("");
    const [filteredResult, setFilteredResult] = useState([]);

    const axiosPrivate = useAxiosPrivate();

    useEffect(() => {
        if (search.trim() === "") {
            setFilteredResult(movies);
            setMessage("");
        } else {
            const filtered = movies.filter(movie =>
                movie.title.toLowerCase().includes(search.toLowerCase())
            );

            setFilteredResult(filtered);

            if (filtered.length === 0) {
                setMessage("No Movies found, Please search again!");
            } else {
                setMessage("");
            }
        }
    }, [search, movies]);

    useEffect(()=>{
        const fetchRecommendedMovies = async() => {
            setLoading(true);
            setMessage("");

            try {
                const resp = await axiosPrivate.get('/recommendedmovies');
                setMovies(resp.data);
                setFilteredResult(resp.data);
                if(resp.data.length === 0){
                    setMessage("No Movies found!!")
                }
            } catch (error) {
                console.log(error);
                setMessage("Something went wrong!!");
            } finally{
                setLoading(false);
            }
        }
        fetchRecommendedMovies();
    }, [axiosPrivate])

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

            {loading ? <h2><Spinner /></h2> : 
            <Movies movies={filteredResult} message={message}/>
            }
        </>
    )

}

export default Recommended;