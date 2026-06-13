import { useState } from "react";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import axiosConfig from "../../api/axiosConfig";
import { Link, useLocation, useNavigate } from "react-router-dom";
import UseAuth from "../../hook/UseAuth";
import logo from "../../../logo.png";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const [message, setMessage] = useState(null);
  const [loading, setLoading] = useState(false);
  const [highlightForgot, setHighlightForgot] = useState(false);

  const { setAuth } = UseAuth();

  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || "/";

  async function handleSubmit(e) {
    e.preventDefault();
    setMessage("");
    setError(null);
    setLoading(true);

    try {
      const loginPayLoad = {
        email: email,
        password: password,
      };

      const res = await axiosConfig.post("/login", loginPayLoad);

      if (res.data.error) {
        setError(res.data.error);
        return;
      }

      setAuth(res.data);
      navigate(from, { replace: true });
    } catch (err) {
      console.log(err);

      const errorMsg = err.response?.data?.error;

      if (errorMsg === "Please verify your email first") {
        const msg = (
          <>
            Please{" "}
            <Button
              className="btn btn-warning"
              onClick={(e) => handleVerifyonLogin(e, err.response.data.user_id)}
            >
              Verify
            </Button>{" "}
            your email first
          </>
        );
        setError(msg);
        setHighlightForgot(false);
        return;
      }

      if (errorMsg === "Invalid credentials") {
        setHighlightForgot(true);
      } else {
        setHighlightForgot(false);
      }

      setError(errorMsg || "Something went wrong, try again!");
    } finally {
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
        user_id: userId,
      });

      navigate(`/verify-email/${userId}`);
    } catch (err) {
      setError(err?.response?.data?.error || "Failed to send OTP");
      console.log(err.response.data.error);
    } finally {
      setLoading(false);
    }
  }

  async function handleForgotPassword(e) {
    e.preventDefault();

    setError("");
    setMessage("");
    setLoading(false);

    if (email.length === 0) {
      setMessage("Enter Email address to Reset Password");
      document.getElementById("email")?.focus();
      return;
    }

    try {
      setLoading(true);
      const res = await axiosConfig.post("/forgot-password", {
        emailid: email,
      });
      setMessage(res.data.message);
    } catch (error) {
      setError(error.response.data.error || "Something went wrong, try again.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Container className="login-container d-flex align-items-center justify-content-center">
      <div className="login-card">
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="login-title">Sign In</h2>
          <p className="login-subtitle">
            Signin to your Magic Movie Stream account.
          </p>
        </div>

        {error && <div className="login-alert login-alert-error">{error}</div>}
        {message && (
          <div className="login-alert login-alert-success">{message}</div>
        )}

        <Form onSubmit={handleSubmit}>
          <Form.Group className="login-field mb-3">
            <Form.Label htmlFor="email">Email</Form.Label>
            <Form.Control
              id="email"
              name="email"
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="login-input"
              required
            />
          </Form.Group>
          <Form.Group className="login-field mb-1">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="login-input"
              required
            />
          </Form.Group>

          <div className="text-end mb-3">
            <Link
              onClick={handleForgotPassword}
              className={`login-forgot ${highlightForgot ? "login-forgot-active" : ""}`}
              style={{
                transition: "all 0.2s ease",
                animation: highlightForgot ? "pulse 1.2s infinite" : "none",
              }}
            >
              Forgot Password?
            </Link>
          </div>

          <Button type="submit" className="login-btn w-100" disabled={loading}>
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                Logging in...
              </>
            ) : (
              "Login"
            )}
          </Button>
        </Form>

        <div className="login-divider">
          <hr />
          <span>or</span>
          <hr />
        </div>

        <div className="login-register-row">
          Don't have an account?{" "}
          {!loading ? <Link to="/register">Resgister here</Link> : <p>Resgister here</p>}
        </div>
      </div>
    </Container>
  );
};

export default Login;
