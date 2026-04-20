import axios from "axios";
import { useEffect } from "react";
import useAuth from "./useAuth";

const apiUrl = import.meta.env.VITE_API_BASE_URL;

// Single shared instance
const axiosAuth = axios.create({
  baseURL: apiUrl,
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});

let isRefreshing = false;
let failedQueue = [];

const useAxiosPrivate = () => {
  const { setAuth } = useAuth();

  //useMemo and useRef are private to components, they are created per component but we want one global, so moving outside.
  // const isRefreshing = useRef(false);
  // const failedQueue = useRef([]);

  const processQueue = (error) => {
    failedQueue.forEach(({ resolve, reject }) => {
      error ? reject(error) : resolve();
    });
    failedQueue = [];
  };

  useEffect(() => {
    const interceptor = axiosAuth.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error?.config;

        if (!originalRequest) {
          return Promise.reject(error);
        }

        if (
          error.response?.status === 401 &&
          !originalRequest._retry &&
          !originalRequest.url.includes("/refresh")
        ) {
          originalRequest._retry = true;

          // If refresh already in progress → queue request to avoid multiple db refresh calls
          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              failedQueue.push({ resolve, reject });
            }).then(() => axiosAuth(originalRequest));
          }

          isRefreshing = true;

          try {
            await axiosAuth.get("/refresh");

            processQueue(null);

            return axiosAuth(originalRequest);
          } catch (err) {
            processQueue(err);

            setAuth(null);

            return Promise.reject(err);
          } finally {
            isRefreshing = false;
          }
        }

        return Promise.reject(error);
      }
    );

    return () => {
      axiosAuth.interceptors.response.eject(interceptor);
    };
  }, [setAuth]);

  return axiosAuth;
};

export default useAxiosPrivate;