import { useState, useEffect } from "react";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import axiosConfig from "../../api/axiosConfig";
import { Link, useNavigate } from "react-router-dom";
import logo from "../../../logo.png";
//  /MagicStreamLogo.png';

const Register = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [adminPassword, setAdminPassword] = useState("");
  const [isAdmin, setIsAdmin] = useState(false);
  const [favouriteGenres, setFavouriteGenres] = useState([]);
  const [genres, setGenres] = useState([]);
  const [error, setError] = useState(null);
  const [message, setMessage] = useState(null);
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  function handleGenreChange(e) {
    const opts = Array.from(e.target.selectedOptions);
    // console.log(opts);
    setFavouriteGenres(
      opts.map((opt) => ({
        genre_id: Number(opt.value),
        genre_name: opt.label,
      })),
    );
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError(null);
    const role = isAdmin ? "ADMIN" : "USER";

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 6) {
      setError("Password must have at least 6 characters");
      document.getElementById("password").focus();
      return;
    }

    setLoading(true);

    try {
      const registerPayLoad = {
        first_name: firstName,
        last_name: lastName,
        email: email,
        role: role,
        password: password,
        adminPassword: adminPassword,
        favourite_genres: favouriteGenres,
      };
      const res = await axiosConfig.post("/register", registerPayLoad);
      if (res.data.error) {
        setError(res.data.error);
        return;
      }
      window.scrollTo({
        top: 0,
        behavior: "smooth",
      });
      setMessage(res?.data?.message);
      setTimeout(() => {
        navigate(`/verify-email/${res.data.user_id}`);
      }, 3000);
    } catch (err) {
      console.log(err);
      setError(
        err.response?.data?.error || "Registration failed, please try again.",
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const FetchGenres = async () => {
      try {
        const genres = await axiosConfig.get("/genres");
        const genredata = genres.data;
        setGenres(genredata);
      } catch (err) {
        console.log(err);
      }
    };
    FetchGenres();
  }, []);

  return (
    <Container className="login-container d-flex align-items-center justify-content-center">
      <div className="login-card" style={{ maxWidth: 400 }}>
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="login-title">Create Account</h2>
          <p className="login-subtitle">Join Magic Stream today</p>
        </div>

        {error && <div className="login-alert login-alert-error">{error}</div>}
        {message && (
          <div className="login-alert login-alert-success">{message}</div>
        )}

        <Form onSubmit={handleSubmit}>
          <div className="row g-3 mb-3">
            <div className="col-6">
              <Form.Group className="login-field">
                <Form.Label>First Name</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter first name"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="login-input"
                  required
                />
              </Form.Group>
            </div>
            <div className="col-6">
              <Form.Group className="login-field">
                <Form.Label>Last Name</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter last name"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="login-input"
                  required
                />
              </Form.Group>
            </div>
          </div>

          <Form.Group className="login-field mb-3">
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="login-input"
              required
            />
          </Form.Group>

          <div className="row g-3 mb-3">
            <div className="col-6">
              <Form.Group className="login-field">
                <Form.Label htmlFor="password">Password</Form.Label>
                <Form.Control
                  id="password"
                  type="password"
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="login-input"
                  required
                />
              </Form.Group>
            </div>
            <div className="col-6">
              <Form.Group className="login-field">
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
            </div>
          </div>

          <div className="row g-3 mb-3">
            <div className="col-6">
              <Form.Group className="login-field">
                <Form.Label>Register as an Admin ?</Form.Label>
                <Form.Select
                  className="login-input"
                  name="isAdmin"
                  value={isAdmin ? "yes" : "no"}
                  onChange={(e) => {
                    setIsAdmin(e.target.value === "yes");
                    if (e.target.value === "no") {
                      setAdminPassword("");
                    }
                  }}
                >
                  <option value="yes" className="admin-option">
                    Yes
                  </option>
                  <option value="no" className="admin-option">
                    No
                  </option>
                </Form.Select>
              </Form.Group>
            </div>
            {isAdmin && (
              <div className="col-6">
                <Form.Group className="login-field">
                  <Form.Label>Admin Password</Form.Label>
                  <Form.Control
                    type="password"
                    placeholder="Adim Password"
                    value={adminPassword}
                    onChange={(e) => setAdminPassword(e.target.value)}
                    className="login-input"
                    required
                  />
                </Form.Group>
              </div>
            )}
          </div>

          <Form.Group className="login-field mb-4">
            <Form.Label>Favourite Genres</Form.Label>
            <Form.Select
              multiple
              value={favouriteGenres.map((g) => String(g.genre_id))}
              onChange={handleGenreChange}
              className="login-input genre-select"
            >
              {genres.map((genre) => (
                <option
                  key={genre.genre_id}
                  value={genre.genre_id}
                  label={genre.genre_name}
                >
                  {genre.genre_name}
                </option>
              ))}
            </Form.Select>
            <Form.Text className="genre-hint">
              Hold Ctrl / Cmd to select multiple genres.
            </Form.Text>
          </Form.Group>

          <Button type="submit" className="login-btn w-100" disabled={loading}>
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                Registering...
              </>
            ) : (
              "Register"
            )}
          </Button>
        </Form>

        <div className="login-divider">
          <hr />
          <span>or</span>
          <hr />
        </div>

        <div className="login-register-row">
          Already have an account?{" "}
          {!loading ? <Link to="/login">login here</Link> : <p>login here</p>}
        </div>
      </div>
    </Container>
  );
};

export default Register;
