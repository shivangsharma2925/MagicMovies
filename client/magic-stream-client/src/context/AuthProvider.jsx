import { createContext, useCallback, useEffect, useMemo, useState } from "react";
import api from "../api/axiosConfig";
import useAxiosPrivate from "../hook/UseAxiosPrivate";

const AuthContext = createContext({});

export const AuthProvider = ({ children }) => {
  const [auth, setAuth] = useState(null);
  const [loading, setLoading] = useState(true);
  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await axiosPrivate.get("/profile/me");
        if (response.data.error) {
          setAuth(null);
          return;
        }
        setAuth(response.data);
      } catch (err) {
        console.log(err);
        setAuth(null);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [axiosPrivate]);

  const logout = useCallback(async () => {
    try {
      await api.post("/logout", { user_id: auth?.user_id });
    } catch (err) {
      console.error("Logout error:", err);
    } finally {
      setAuth(null);
    }
  }, [auth]);

  const value = useMemo(() => {
    return { auth, setAuth, logout, loading };
  }, [auth, loading, logout]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export default AuthContext;
