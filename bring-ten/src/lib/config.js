export const apiURL = process.env.VITE_API_URL ?? import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export const publicApiURL = import.meta.env.VITE_PUBLIC_API_URL ?? 'http://localhost:8080';
