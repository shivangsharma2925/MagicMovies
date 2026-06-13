import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { useNavigate, NavLink, Link } from "react-router-dom";
import { useEffect, useRef } from "react";
import UseAuth from "../../hook/UseAuth";
import logo from "../../../logo.png";

const Header = () => {
  const navigate = useNavigate();
  const { auth, logout } = UseAuth();

  const initials = auth ? `${auth.first_name?.[0] ?? ""}${auth.last_name?.[0] ?? ""}`.toUpperCase()
    : "";

  const navRef = useRef(null);

  useEffect(() => {
    const el = navRef.current;
    if (!el) return;

    const updateHeight = () => {
      document.documentElement.style.setProperty("--navbar-h", `${el.offsetHeight}px`);
    };

    updateHeight(); // initial

    // Watches for any DOM/class changes inside the navbar (collapse open/close)
    const observer = new ResizeObserver(updateHeight);
    observer.observe(el);

    return () => observer.disconnect();
  }, []);

  return (
    <Navbar ref={navRef} expand="lg" sticky="top" className="movie-navbar">
      <Container>
        <Navbar.Brand as={NavLink} to="/" className="d-flex align-items-center brand-glow gap-2">
          <img
            alt=""
            src={logo}
            width="34"
            height="34"
            className="me-2 rounded"
          />
          <span className="brand-text">Magic Stream</span>
        </Navbar.Brand>

        <Navbar.Toggle aria-controls="main-navbar-nav" className="navbar-toggler-custom" />

        <Navbar.Collapse id="main-navbar-nav">
          <Nav className="me-auto gap-1 mt-2 mt-lg-0">
            <Nav.Link as={NavLink} to="/" className="nav-item-custom">
              Home
            </Nav.Link>
            <Nav.Link
              as={NavLink}
              to="/recommended"
              className="nav-item-custom"
            >
              Recommended
            </Nav.Link>
          </Nav>

          <Nav className="ms-auto align-items-center gap-2 mt-lg-0 mb-2 mb-lg-0">
            {auth ? (
              <>
                <div className="user-chip">
                  <div className="user-avatar">{initials}</div>
                  <span className="user-name-text">{auth.first_name}</span>
                </div>
                <Button className="btn-logout" size="sm" onClick={logout}>
                  Logout
                </Button>
              </>
            ) : (
              <>
                <div className="d-flex gap-2">
                  <Button
                    className="btn-login"
                    size="sm"
                    onClick={() => navigate("/login")}
                  >
                    Login
                  </Button>

                  <Button
                    className="btn-register"
                    size="sm"
                    onClick={() => navigate("/register")}
                  >
                    Register
                  </Button>
                </div>
              </>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Header;
