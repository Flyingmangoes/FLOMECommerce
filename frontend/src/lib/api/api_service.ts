export interface FetchOptions {
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
    url: string;
    headers?: Record<string, string>;
    body?: unknown;
}

export interface FetchResponse<T = unknown> {
    code: number;
    message: string;
    data: T | null;
}

const ALLOWED_METHODS = ['GET', 'POST', 'PUT', 'DELETE'] as const;

export class ApiService {
    private readonly serverUrl: string;
    private readonly apiVersion: string = 'v1';

    constructor(serverUrl: string, apiVersion?: string) {
        this.serverUrl = serverUrl;
        if (apiVersion) {
            this.apiVersion = apiVersion;
        }
    }

    async checkConnection(): Promise<boolean> {
        if (!this.serverUrl?.trim()) {
            console.error('Server URL is not set');
            return false;
        }

        try {
            const res = await fetch(`${this.serverUrl}/${this.apiVersion}/health`);
            return res.ok;
        } catch (error) {
            console.error('Error checking API connection:', error);
            return false;
        }
    }

    async fetch<T = unknown>(options: FetchOptions): Promise<FetchResponse<T>> {
        const method = options.method ?? 'GET';

        if (!options.url) {
            return { code: 400, message: 'URL is required', data: null };
        }
        if (!ALLOWED_METHODS.includes(method)) {
            return { code: 400, message: 'Method is not allowed', data: null };
        }

        try {
            const response = await fetch(`${this.serverUrl}/${this.apiVersion}${options.url}`, {
                method,
                headers: {
                'Content-Type': 'application/json',
                ...options.headers,
                },
                body: options.body ? JSON.stringify(options.body) : undefined,
            });

            const text = await response.text();
            let data: T | null = null;
            if (text) {
                try {
                    data = JSON.parse(text) as T;
                } catch {
                    data = text as unknown as T;
                }
            }

            return { code: response.status, message: response.statusText, data };
        } catch (error) {
            console.error('Error in ApiService.fetch:', error);
            return { code: 0, message: 'Network error \u2014 could not reach server', data: null };
        }
    }
}