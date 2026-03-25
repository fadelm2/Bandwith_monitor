import axios from 'axios';

const client = axios.create({
    baseURL: 'http://localhost:9090/internal',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Attach token to requests
client.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Handle 401 Unauthorized globally
client.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('username');
            window.location.reload();
        }
        return Promise.reject(error);
    }
);

export const api = {
    health: () => client.get('/health'),

    // Capacity CRUD
    listCapacity: () => client.get('/capacity'),
    getCapacity: (id) => client.get(`/capacity/${id}`),
    createCapacity: (data) => client.post('/capacity', data),
    updateCapacity: (id, data) => client.put(`/capacity/${id}`, data),
    deleteCapacity: (id) => client.delete(`/capacity/${id}`),
    bulkUpdateCapacity: (data) => client.put('/capacity/bulk', data),

    // Traffic
    searchTraffic: (params) => client.get('/traffic', { params }),
};

export default client;
