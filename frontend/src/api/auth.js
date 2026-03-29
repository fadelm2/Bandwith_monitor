import axios from 'axios';

// Separate auth client pointing to /api
const authClient = axios.create({
    baseURL: import.meta.env.VITE_AUTH_API_URL || 'http://localhost:9090/api',
    headers: { 'Content-Type': 'application/json' },
});

// Attach token to auth requests
authClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
});

// Handle 401 Unauthorized globally for auth client too
authClient.interceptors.response.use(
    (response) => response,
    (error) => {
        // Only trigger on actual 401, but NOT on the login call itself 
        // (login 401 is handled locally in LoginPage.jsx)
        if (error.response && error.response.status === 401 && !error.config.url.includes('/auth/login')) {
            localStorage.removeItem('token');
            localStorage.removeItem('username');
            window.location.reload();
        }
        return Promise.reject(error);
    }
);

export const authApi = {
    login: (id, password) =>
        authClient.post('/auth/login', { id, password }),

    logout: () =>
        authClient.post('/auth/logout'),

    current: () =>
        authClient.get('/auth/current'),
};
