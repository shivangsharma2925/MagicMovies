import { useState } from "react";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import Container from "react-bootstrap/Container";
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

  const handleResetPassword = async (e) => {
    e.preventDefault();

    setError("");
    setLoading(false);
    try {
      setLoading(true);
      const res = await axiosConfig.post("/reset-password", {
        token: uniqueKey,
        password: password,
      });
      setMessage(res.data.message);
      setTimeout(() => {
        // navigate("/login");
        window.close();
      }, 1500);
    } catch (err) {
      setError(
        err.response.data.error || "Unable to Reset password, try again.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="login-container d-flex align-items-center justify-content-center">
      <div className="login-card">
        <div className="text-center mb-4">
          <div className="login-logo fs-1">🔒</div>
          <h2 className="login-title">Reset Password</h2>
          <p className="login-subtitle">
            Choose a new password for your account
          </p>
        </div>

        {error && <div className="login-alert login-alert-error">{error}</div>}
        {message && (
          <div className="login-alert login-alert-success">{message}</div>
        )}

        <Form onSubmit={handleResetPassword}>
          <Form.Group className="login-field mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter new password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="login-input"
              required
            />
          </Form.Group>

          <Form.Group className="login-field mb-4">
            <Form.Label>Confirm Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Confirm Password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className="login-input"
              isInvalid={!!confirmPassword && password !== confirmPassword}
            />
            <Form.Control.Feedback type="invalid" style={{ fontSize: 11 }}>
              Passwords do not match.
            </Form.Control.Feedback>
          </Form.Group>
          <Button className="login-btn w-100" type="submit" disabled={loading}>
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                'Updating...'
              </>
            ) : (
              "Update Password"
            )}
          </Button>
        </Form>
      </div>
    </Container>
  );
};

export default ResetPassword;
