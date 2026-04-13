import { useEffect, useState } from "react";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import Movies from "../movies/Movies";

const Recommended = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();

    const axiosPrivate = useAxiosPrivate();

    useEffect(()=>{
        const fetchRecommendedMovies = async() => {
            setLoading(true);
            setMessage("");

            try {
                const resp = await axiosPrivate.get('/recommendedmovies');
                setMovies(resp.data);
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
            {loading ? <h2>loading...</h2> : 
            <Movies movies={movies} message={message}/>
            }
        </>
    )

}

export default Recommended;