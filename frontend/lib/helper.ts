export function getBaseURL() {
    const isServer = typeof window === 'undefined';
  
    if (isServer) {
      // Server-side (inside Docker container)
      if (process.platform === 'darwin') {
        // macOS Docker
        return 'http://docker.for.mac.localhost:8080';
      } else if (process.platform === 'win32') {
        // Windows Docker
        return 'http://host.docker.internal:8080';
      } else {
        // Linux Docker
        return 'http://host.docker.internal:8080'; // Default Docker bridge IP
      }
    } else {
      // Client-side (Browser)
      return 'http://localhost:8080'; // Use relative URLs for client-side
    }
  }