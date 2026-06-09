import { useState } from "react";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosConfig from "../../api/axiosConfig";
import { FormGroup } from "react-bootstrap";
// import { useNavigate } from "react-router-dom";

const ResetPassword = () => {
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [message, setMessage] = useState("");

    // const navigate = useNavigate();
    const uniqueKey = new URLSearchParams(location.search).get("uniqueKey");
    
    const handleResetPassword = async(e) => {
        e.preventDefault();

        setError("");
        setLoading(false);
        try{
            setLoading(true);
            const res = await axiosConfig.post("/reset-password", {
                "token": uniqueKey,
                "password": password
            })
            setMessage(res.data.message);
            setTimeout(() => {
                // navigate("/login");
                window.close();
            }, 1500);
        }catch(err){
            setError(err.response.data.error || "Unable to Reset password, try again.")
        }finally{
            setLoading(false);
        }
    }

    return(
        <div className="d-flex justify-content-center align-items-center" style={{"margin": "10%", "flex-direction": "column"}}>
            {error && <div className="alert alert-danger py-2">{error}</div>}
            {message && <div className="alert alert-success py-2">{message}</div>}
            <Form onSubmit={handleResetPassword}>
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
                <Button
                    variant="primary"
                    type="submit"
                    className="w-100 mb-2"
                    disabled={loading}
                    style={{fontWeight: 600, letterSpacing: 1}}
                >
                    {loading ? (
                    'Updating...'
                    ) : 'Update Password'}
                </Button>
            </Form>
        </div>
    )
}

export default ResetPassword;