import { SERVICE_API_HOST } from '$env/static/private';

export async function load({ locals }) {
  const user = locals.user;
    
  try {
      const [shiftsResponse, requestsResponse] = await Promise.all([
          fetch(`http://${SERVICE_API_HOST}/api/shift/assigned`, {
              headers: {
                  'Authorization': `Bearer ${locals.token}`
              }
          }),
          fetch(`http://${SERVICE_API_HOST}/api/shift-request/pending`, {
              headers: {
                  'Authorization': `Bearer ${locals.token}`
              }
          })
      ]);

      // If unauthorized, clear the token
      if (shiftsResponse.status === 401 || requestsResponse.status === 401) {
          cookies.delete('token', { path: '/' });
          cookies.delete('user', { path: '/' });
          return {
              user,
              shifts: [],
              shiftsRequest: [],
              status: 401,
              error: 'Session expired. Please login again.'
          };
      }

      const [shiftsData, requestsData] = await Promise.all([
          shiftsResponse.json(),
          requestsResponse.json()
      ]);
      
      return {
          user,
          shifts: shiftsData.data ?? [],
          shiftsRequest: requestsData.data ?? [],
          status: shiftsResponse.status,
          error: shiftsData.status !== 200 ? (shiftsData.message ?? 'Unknown error.') : null
      };
  } catch (error) {
      console.error('Error fetching data:', error);
      return {
          user,
          shifts: [],
          shiftsRequest: [],
          status: 500,
          error: 'Failed to fetch data'
      };
  }
}