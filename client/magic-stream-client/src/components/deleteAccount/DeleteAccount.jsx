import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import logo from "../../assets/logo.png";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import axiosConfig from "../../api/axiosConfigwocredential";
import {
  faEye,
  faEyeSlash,
  faTriangleExclamation,
  faXmark,
} from "@fortawesome/free-solid-svg-icons";
import UseAuth from "../../hook/UseAuth";
import PageTitle from "../../utils/PageTitle";
import "../../styling/LoginRegister.css";

const DeleteAccount = () => {
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [error, setError] = useState(null);
  const [message, setMessage] = useState(null);
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const { auth, setAuth } = UseAuth();
  const axiosPrivate = useAxiosPrivate();

  const handleOpenConfirmation = (e) => {
    e.preventDefault();
    if (!password.trim()) {
      setError("Please enter your password.");
      return;
    }
    setError("");
    setShowConfirmation(true);
  };

  const handleDeleteAccount = async (e) => {
    e.preventDefault();

    setError("");
    setMessage("");
    setLoading(false);

    try {
      setLoading(true);
      const res = await axiosPrivate.post("/delete-account", {
        password: password,
      });
      setAuth(null);
      setMessage(res.data.message);
      navigate("/login");
    } catch (error) {
      setError(error.response.data.error || "Something went wrong, try again.");
      setShowConfirmation(false);
    } finally {
      setLoading(false);
    }
  };

  async function handleForgotPassword(e) {
    e.preventDefault();

    setError("");
    setMessage("");
    setLoading(false);

    try {
      setLoading(true);
      const res = await axiosConfig.post("/forgot-password", {
        emailid: auth.email,
      });
      setMessage(res.data.message);
    } catch (error) {
      setError(error.response.data.error || "Something went wrong, try again.");
    } finally {
      setLoading(false);
    }
  }

  const handleCloseConfirmation = () => {
    setShowConfirmation(false);
  };

  return (
    <>
      <PageTitle title="Delete Account"/>
      <Container className="delete-account-section d-flex align-items-center justify-content-center">
        <div className="delete-account-card">
          <div className="text-center mb-4">
            <img src={logo} alt="Logo" width={60} className="mb-2" />
            <h2 className="delete-account-title">Delete account</h2>
            <p className="delete-account-subtitle">
              Permanently delete your account and all data associated with it.
            </p>
          </div>

          {error && (
            <div className="delete-account-alert delete-account-alert-error">
              {error}
            </div>
          )}
          {message && (
            <div className="delete-account-alert delete-account-alert-success">
              {message}
            </div>
          )}

          <Form onSubmit={handleOpenConfirmation}>
            <Form.Group className="delete-account-field mb-1">
              <Form.Label>Confirm your password</Form.Label>
              <div className="delete-account-password-wrapper">
                <Form.Control
                  type={showPassword ? "text" : "password"}
                  placeholder="Enter your current password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="delete-account-input"
                  required
                />

                <button
                  type="button"
                  className="delete-password-toggle"
                  onClick={() => setShowPassword((previous) => !previous)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  <FontAwesomeIcon icon={showPassword ? faEyeSlash : faEye} />
                </button>
              </div>
            </Form.Group>

            <div className="text-end mb-4">
              <Link
                onClick={handleForgotPassword}
                className="delete-account-forgot"
              >
                Forgot Password?
              </Link>
            </div>

            <Button
              type="submit"
              className="delete-account-button w-100"
              disabled={loading}
            >
              {loading ? (
                <>
                  <span
                    className="spinner-border spinner-border-sm me-2"
                    role="status"
                    aria-hidden="true"
                  ></span>
                  Sending Link...
                </>
              ) : (
                "Delete Account"
              )}
            </Button>
          </Form>
        </div>
      </Container>

      {showConfirmation && (
        <div
          className="delete-confirm-overlay"
          onClick={handleCloseConfirmation}
        >
          <div
            className="delete-confirm-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="delete-confirm-title"
            onClick={(e) => e.stopPropagation()}
          >
            <button
              type="button"
              className="delete-confirm-close"
              onClick={handleCloseConfirmation}
              disabled={loading}
              aria-label="Close confirmation"
            >
              <FontAwesomeIcon icon={faXmark} />
            </button>

            <div className="delete-confirm-icon">
              <FontAwesomeIcon icon={faTriangleExclamation} />
            </div>

            <h2 id="delete-confirm-title">Permanently delete your account?</h2>

            <p>
              This action cannot be undone. Your account, watchlists, saved
              movies, preferences, and all other associated data will be
              permanently deleted.
            </p>

            <div className="delete-confirm-warning">
              This action is non-recoverable.
            </div>

            <div className="delete-confirm-actions">
              <button
                type="button"
                className="delete-cancel-button"
                onClick={handleCloseConfirmation}
                disabled={loading}
              >
                Cancel
              </button>

              <button
                type="button"
                className="delete-confirm-button"
                onClick={(e) => handleDeleteAccount(e)}
                disabled={loading}
              >
                {loading ? "Deleting..." : "Delete permanently"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default DeleteAccount;
