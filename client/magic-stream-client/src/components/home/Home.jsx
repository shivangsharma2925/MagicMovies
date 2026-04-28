import { useEffect, useState } from "react";
import axiosConfig from "../../api/axiosConfig";
import Movies from "../movies/Movies";
import { Button, Form } from "react-bootstrap";
import Spinner from "../../utils/Spinner";
import { useNavigate } from "react-router-dom";

const Home = ()=>{
    const [movies, setMovies] = useState([])
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState()
    const [search, setSearch] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState(""); 
    const navigate = useNavigate();

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

                if(response.data == null || response.data.length === 0){
                    setMessage(`There are currently no movies available ${trimmedSearch?"for the provided search":""}`);
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
            <div 
                className="sticky-top py-3" 
                style={{ 
                    zIndex: 1020, 
                    top: "60px",
                    backgroundColor: "rgba(255, 255, 255, 0.7)", // Semi-transparent white
                    backdropFilter: "blur(10px)",                // Blurs content behind
                    WebkitBackdropFilter: "blur(10px)",          // Safari support
                }}
            >
                <div className="container-fluid d-flex align-items-center">
                    <div style={{ flex: 1 }}></div>
                    
                    <div style={{ flex: 2 }}>
                        <Form.Control
                            className="shadow-sm border-0 bg-light" 
                            type="text"
                            placeholder="Search movies..."
                            style={{ borderRadius: "20px" }}     
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                        />
                    </div>

                    <div style={{ flex: 1, textAlign: 'right' }}>
                        <Button variant="primary" size="sm" className="rounded-pill px-2 shadow-sm" onClick={() => navigate("admin/add-movie")}>
                            Add Movies
                        </Button>
                    </div>
                </div>
            </div>
            {loading ? 
                (<h2><Spinner /></h2>)
            : (<Movies movies={movies} message={message} />)
            }
        </>
    );
}

export default Home;