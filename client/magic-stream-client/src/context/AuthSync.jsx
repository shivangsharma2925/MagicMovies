import { useEffect } from "react";
import { useLocation } from "react-router-dom";
import UseAuth from "../hook/UseAuth";

const AuthSync = () => {
  const { setAuth } = UseAuth();
  const location = useLocation();

  useEffect(() => {
    const user = localStorage.getItem("user");

    if (!user) {
      setAuth(null);
    }
  }, [location, setAuth]); // runs on every route change

  return null;
};

export default AuthSync;