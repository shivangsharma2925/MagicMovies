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
    const [message, setMessage] = useState(null);
    const [loading, setLoading] = useState(false);
    const {setAuth} = UseAuth();

    const navigate = useNavigate()
    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    async function handleSubmit(e) {
        e.preventDefault();
        setMessage("");
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
            navigate(from, {replace: true});
        } catch (err) {
            console.log(err);
            if(err.response?.data?.error && err.response?.data?.error === "Please verify your email first"){
                const msg = <>Please <Button onClick={(e)=>handleVerifyonLogin(e, err.response.data.user_id)}>Verify</Button> your email first</>;
                setError(msg);
            }else setError(err.response?.data?.error || "Something went wrong, try again!")
        } finally{
            setLoading(false);
        }
    }

    async function handleVerifyonLogin(event, userId) {
        event.preventDefault();

        setError("");
        setMessage("");
        setLoading(true);
        try {
            await axiosConfig.post("/resend-verification", {
                "user_id": userId
            })

            navigate(`/verify-email/${userId}`);

        } catch (err) {
            setError(err?.response?.data?.error || "Failed to send OTP");
            console.log(err.response.data.error);
        } finally{
            setLoading(false);
        }
    } 

    async function handleForgotPassword(e){
        e.preventDefault();

        setError("");
        setMessage("");
        setLoading(false);

        if(email.length === 0){
            setError("Enter Email address to Reset Password");
            return;
        }

        try {
            setLoading(true);
            const res =await axiosConfig.post("/forgot-password", {
                "emailid": email
            })
            setMessage(res.data.message);
        } catch (error) {
            setError(error.response.data.error || "Something went wrong, try again.")
        }finally{
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
                    {message && <div className="alert alert-success py-2">{message}</div>}                
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
                        <div className="text-end">
                            <small>
                                <Link onClick={(e)=>handleForgotPassword(e)} className="text-muted text-decoration-none">Forgot Password?</Link>
                            </small>
                        </div>
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