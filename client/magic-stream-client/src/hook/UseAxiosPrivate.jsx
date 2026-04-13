import axios from "axios";
import { useMemo, useRef, useEffect } from "react";
import useAuth from "./UseAuth";

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {
  const { setAuth } = useAuth();

  const isRefreshing = useRef(false);
  const failedQueue = useRef([]);

  const axiosAuth = useMemo(() => {
    return axios.create({
      baseURL: apiUrl,
      headers: { "Content-Type": "application/json" },
      withCredentials: true,
    });
  }, []);

  const processQueue = (error) => {
    failedQueue.current.forEach(({ resolve, reject }) => {
      if (error) reject(error);
      else resolve();
    });
    failedQueue.current = [];
  };

  useEffect(() => {
    const interceptor = axiosAuth.interceptors.response.use(
      (res) => res,
      async (error) => {
        const originalRequest = error.config;

        if (
          error.response?.status === 401 &&
          !originalRequest._retry
        ) {
          originalRequest._retry = true;

          if (isRefreshing.current) {
            return new Promise((resolve, reject) => {
              failedQueue.current.push({ resolve, reject });
            }).then(() => axiosAuth(originalRequest));
          }

          isRefreshing.current = true;

          try {
            await axiosAuth.post("/refresh");
            processQueue(null);
            return axiosAuth(originalRequest);
          } catch (err) {
            processQueue(err);
            localStorage.removeItem('user');
            setAuth(null);
            return Promise.reject(err);
          } finally {
            isRefreshing.current = false;
          }
        }

        return Promise.reject(error);
      }
    );

    return () => axiosAuth.interceptors.response.eject(interceptor);
  }, [axiosAuth, setAuth]);

  return axiosAuth;
};

export default useAxiosPrivate;