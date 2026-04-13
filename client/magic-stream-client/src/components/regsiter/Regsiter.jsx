import { useState, useEffect } from "react";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import axiosConfig from "../../api/axiosConfig";
import { Link, Navigate } from "react-router-dom";
import logo from '../../assets/MagicStreamLogo.png';

const Register = () => {
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [favouriteGenres, setFavouriteGenres] = useState([]);
    const [genres, setGenres] = useState([]);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    function handleGenreChange(e){
        const opts = Array.from(e.target.selectedOptions);
        // console.log(opts);
        setFavouriteGenres(opts.map(opt => (
            {
                genre_id: Number(opt.value),
                genre_name: opt.label
            }
        )));
    }

    async function handleSubmit(e){
        e.preventDefault();
        setError(null);
        const defaultRole = "USER"

        if(password !== confirmPassword){
            setError("Passwords do not match");
            return;
        }

        setLoading(true);

        try{
            const registerPayLoad = {
                "first_name": firstName,
                "last_name": lastName,
                "email": email,
                "role": defaultRole,
                "password": password,
                "favourite_genres": favouriteGenres
            };
            const res = await axiosConfig.post("/register", registerPayLoad);
            if(res.data.error){
                setError(res.data.error);
                return;
            }
            <Navigate to='/login' replace />
        }catch(err){
            console.log(err);
            setError("Registration failed, please try again.");
        }finally{
            setLoading(false);
        }
    }

    useEffect(()=>{
        const FetchGenres = async() => {
            try{
                const genres = await axiosConfig.get("/genres");
                const genredata = genres.data;
                setGenres(genredata);
            }catch(err){
                console.log(err);
            }
        }
        FetchGenres();
    }, [])

    return (
        <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
            <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth: 400, width: '100%'}}>
                <div className="text-center mb-4">
                    <img src={logo} alt="Logo" width={60} className="mb-2" />
                    <h2 className="fw-bold">Register</h2>
                    <p className="text-muted">Create your Magic Movie Stream account.</p>
                    {error && <div className="alert alert-danger py-2">{error}</div>}                
                </div>
                <Form onSubmit={handleSubmit}>
                    <Form.Group className="mb-3">
                        <Form.Label>First Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter first name"
                            value={firstName}
                            onChange={e => setFirstName(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Last Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter last name"
                            value={lastName}
                            onChange={e => setLastName(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Email</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={e => setEmail(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={e => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Confirm Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Confirm Password"
                            value={confirmPassword}
                            onChange={e => setConfirmPassword(e.target.value)}
                            required
                            isInvalid ={!!confirmPassword && password !== confirmPassword}

                        />
                        <Form.Control.Feedback type="invalid">
                            Passwords do not match.
                        </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Select
                            multiple
                            value={favouriteGenres.map(g => String(g.genre_id))}
                            onChange={handleGenreChange}
                        >
                            {genres.map(genre => (
                                <option key={genre.genre_id} value={genre.genre_id} label={genre.genre_name}>
                                    {genre.genre_name}
                                </option>
                            ))}
                        </Form.Select>
                        <Form.Text className="text-muted">
                            Hold Ctrl (Windows) or Cmd (Mac) to select multiple genres.
                        </Form.Text>
                    </Form.Group>
                    <Button
                        variant="primary"
                        type="submit"
                        className="w-100 mb-2"
                        disabled={loading}
                        style={{fontWeight: 600, letterSpacing: 1}}
                    >
                        {loading ? (
                            <>
                                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                                Registering...
                            </>
                        ) : 'Register'}
                    </Button>                        
                </Form>
                <div className="text-center mt-3">
                    <span className="text-muted">Already have an account? </span>
                    <Link to="/login" className="fw-semibold">login here</Link>
                </div>
            </div>           
       </Container>
    )
}

export default Register;