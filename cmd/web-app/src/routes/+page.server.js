export async function load({ locals }) {
    return {
      user: JSON.parse(locals.user),
    };
}