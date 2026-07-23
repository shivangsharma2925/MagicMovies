import { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowUp } from "@fortawesome/free-solid-svg-icons";

const ScrollToTopBtn = () => {
  const [show, setShow] = useState(false);
  const [hasScrollbar, setHasScrollbar] = useState(false);

  useEffect(() => {
    const checkScrollbar = () => {
      setHasScrollbar(document.documentElement.scrollHeight > window.innerHeight); //if whole div height is greater than the fixed innner height, then there is a scrollbar
    };

    checkScrollbar();

    window.addEventListener("resize", checkScrollbar);

    return () => window.removeEventListener("resize", checkScrollbar);
  }, []);

  useEffect(() => {
    const onScroll = () => {
      setShow(window.scrollY > 200); //if user has scrolled more than 200px, then display the button
    };

    window.addEventListener("scroll", onScroll);

    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  const handleClick = () => {
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };

  if (!(hasScrollbar && show)) return null;

  return (
    <button style={styles} onClick={handleClick}>
      <FontAwesomeIcon icon={faArrowUp} />
    </button>
  );
};

const styles = {
  position: "fixed",
  right: "20px",
  bottom: "20px",
  width: "42px",
  height: "42px",
  borderRadius: "50%",
  border: "none",
  cursor: "pointer",
  color: "#fff",
  zIndex: 5001,
  background: "linear-gradient(135deg, #ff2e63, #7928ca)",
  boxShadow: "0 0 8px rgba(255,255,255,0.5)",
};

export default ScrollToTopBtn;
