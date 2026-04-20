import { useState } from "react";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import axiosConfig from "../../api/axiosConfig";
import { Link, useLocation, useNavigate } from "react-router-dom";
import UseAuth from "../../hook/UseAuth";
import logo from '../../assets/MagicStreamLogo.png';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const {setAuth} = UseAuth();

    const naviagte = useNavigate()
    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    async function handleSubmit(e) {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            const loginPayLoad = {
                "email": email,
                "password": password
            }
            const res = await axiosConfig.post('/login', loginPayLoad);
            if(res.data.error){
                setError(res.data.error);
                return;
            }
            setAuth(res.data);
            // localStorage.setItem('user', JSON.stringify(res.data));
            naviagte(from, {replace: true});
        } catch (error) {
            console.log(error);
            setError("Something went wrong, try again!")
        } finally{
            setLoading(false);
        }
    }

    return (
        <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
            <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth: 400, width: '100%'}}>
                <div className="text-center mb-4">
                    <img src={logo} alt="Logo" width={60} className="mb-2" />
                    <h2 className="fw-bold">Sign In</h2>
                    <p className="text-muted">Signin to your Magic Movie Stream account.</p>
                    {error && <div className="alert alert-danger py-2">{error}</div>}                
                </div>
                <Form onSubmit={handleSubmit}>
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
                                Logging in...
                            </>
                        ) : 'Login'}
                    </Button>                        
                </Form>
                <div className="text-center mt-3">
                    <span className="text-muted">Don't have an account? </span>
                    <Link to="/register" className="fw-semibold">Register here</Link>
                </div>
            </div>          
       </Container>
    )
}

export default Login;