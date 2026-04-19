import { useEffect, useState } from "react";
import axiosConfig from "../../api/axiosConfig";
import Movies from "../movies/Movies";
import { Form } from "react-bootstrap";
import Spinner from "../../utils/Spinner";

const Home = ({updateMovieReview})=>{
    const [movies, setMovies] = useState([])
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState()
    const [search, setSearch] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState(""); 

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(search);
        }, 500); // 500ms delay

        return () => clearTimeout(timer);
    }, [search]);

    useEffect(()=>{
        const fetchMovies = async ()=>{
            setLoading(true);
            setMessage("");
            try{
                const trimmedSearch = debouncedSearch.trim();

                const url = trimmedSearch
                ? `/movies?search=${trimmedSearch}`
                : `/movies`;

                const response = await axiosConfig.get(url);

                setMovies(response.data);

                if(response.data.length === 0){
                    setMessage("There are currently no movies available");
                }
            }catch(error){
                console.log(error);
                setMessage("Error fetching movies");
            }finally{
                setLoading(false);
            }
        }
        fetchMovies();
    }, [debouncedSearch])

    return (
        <>
            <Form className="mb-3">
                <Form.Control
                    type="text"
                    placeholder="Search movies..."
                    value={search}
                    className="w-50 m-auto mt-2"
                    onChange={(e) => setSearch(e.target.value)}
                />
            </Form>
            {loading ? 
                (<h2><Spinner /></h2>)
            : (<Movies movies={movies} message={message} updateMovieReview={updateMovieReview} />)
            }
        </>
    );
}

export default Home;