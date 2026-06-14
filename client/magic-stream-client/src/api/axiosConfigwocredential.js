import axios from 'axios';

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const apiwocredential = axios.create({
    baseURL: apiUrl,
    headers: {'Content-Type': 'application/json'},
    withCredentials: false
})

export default apiwocredential;